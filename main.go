package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/capitalone/fpe/ff1"
)

// ── Steganography on tweak bytes ──────────────────────────────────────────────
//
// Byte alphabet: 0-255 (256 symbols).
// sft = 4:
//   low  4 bits = random noise (preserved from original random tweak)
//   high 4 bits = hidden data  (keyID or version, each 0-15)
//   max value   = (15 << 4) + 15 = 255 — fits exactly in one byte ✓
//
// Mirrors Ubiq encodeKeyNumber / decodeKeyNumber:
//   encode: idx = posOf(byte); idx += n << sft; byte = valAt(idx)
//   decode: c   = posOf(byte); n = c >> sft;    byte = valAt(c - (n<<sft))
// For raw bytes, posOf and valAt are identity functions.
//
// Tweak layout (8 bytes, all start as random):
//   [0]   keyID   hidden in high nibble, original random in low nibble
//   [1]   version hidden in high nibble, original random in low nibble
//   [2:8] pure random nonce — untouched
//
// Supports keyID 0-15, version 0-15

const sft = 4

// encodeByte hides n (0-15) into the high nibble of random byte r.
// Low nibble of r preserved as noise — matching Ubiq's encodeKeyNumber.
func encodeByte(r byte, n int) byte {
	low := int(r) & 0x0f
	return byte(low + (n << sft))
}

// decodeByte extracts n from the high nibble — matching Ubiq's decodeKeyNumber.
func decodeByte(b byte) int {
	return int(b) >> sft
}

// buildTweak:
//
//	step 1 — generate 8 fully random bytes
//	step 2 — stego: encode keyID into tweak[0], version into tweak[1]
//	step 3 — tweak[2:8] remain purely random and untouched
func buildTweak(keyID, version int) []byte {
	if keyID < 0 || keyID > 15 {
		log.Fatalf("keyID must be 0-15, got %d", keyID)
	}
	if version < 0 || version > 15 {
		log.Fatalf("version must be 0-15, got %d", version)
	}

	tweak := make([]byte, 8)
	if _, err := rand.Read(tweak); err != nil {
		log.Fatalf("rand: %v", err)
	}

	fmt.Printf("  step 1 — random tweak (before stego):\n")
	fmt.Printf("           %x\n", tweak)
	fmt.Printf("           tweak[0] = %08b (random)\n", tweak[0])
	fmt.Printf("           tweak[1] = %08b (random)\n", tweak[1])

	tweak[0] = encodeByte(tweak[0], keyID)
	tweak[1] = encodeByte(tweak[1], version)

	fmt.Printf("\n  step 2 — after stego encode:\n")
	fmt.Printf("           tweak[0] = %08b  hi=%d (keyID=%d)    lo=%d (rand)\n",
		tweak[0], tweak[0]>>sft, keyID, tweak[0]&0x0f)
	fmt.Printf("           tweak[1] = %08b  hi=%d (version=%d)  lo=%d (rand)\n",
		tweak[1], tweak[1]>>sft, version, tweak[1]&0x0f)
	fmt.Printf("           tweak[2:] = %x (pure random, untouched)\n", tweak[2:])
	fmt.Printf("\n  step 3 — final tweak: %x\n", tweak)

	return tweak
}

// extractTweak decodes keyID and version from the tweak without mutating it.
func extractTweak(tweak []byte) (keyID, version int) {
	keyID = decodeByte(tweak[0])
	version = decodeByte(tweak[1])
	return
}

// ── Keys ──────────────────────────────────────────────────────────────────────

// lookupKey derives a deterministic 32-byte AES-256 key from keyID + version.
// Any keyID (0-15) and version (0-15) is valid — no static table needed.
// In production this is a real KMS lookup.
func lookupKey(keyID, version int) []byte {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte((keyID*31 + version*17 + i*7) & 0xff)
	}
	return key
}

// ── FPE ───────────────────────────────────────────────────────────────────────

func fpeEncrypt(plaintext string, keyID, version int) (string, []byte) {
	tweak := buildTweak(keyID, version)
	key := lookupKey(keyID, version)

	c, err := ff1.NewCipher(10, len(tweak), key, tweak)
	if err != nil {
		log.Fatalf("ff1 init: %v", err)
	}
	ct, err := c.Encrypt(plaintext)
	if err != nil {
		log.Fatalf("ff1 encrypt: %v", err)
	}
	return ct, tweak
}

func fpeDecrypt(ciphertext string, tweak []byte) (string, int, int) {
	fmt.Printf("  raw tweak         : %x\n", tweak)
	fmt.Printf("  tweak[0]          : %08b\n", tweak[0])
	fmt.Printf("  tweak[1]          : %08b\n", tweak[1])

	keyID, version := extractTweak(tweak)

	fmt.Printf("  extracted keyID   : tweak[0] >> %d = %d\n", sft, keyID)
	fmt.Printf("  extracted version : tweak[1] >> %d = %d\n", sft, version)

	key := lookupKey(keyID, version)
	c, err := ff1.NewCipher(10, len(tweak), key, tweak)
	if err != nil {
		log.Fatalf("ff1 init: %v", err)
	}
	pt, err := c.Decrypt(ciphertext)
	if err != nil {
		log.Fatalf("ff1 decrypt: %v", err)
	}
	return pt, keyID, version
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	plaintext := "1234567890123456"
	keyID := 15
	version := 15

	fmt.Println("════ ENCRYPT ════════════════════════════════════════")
	fmt.Printf("plaintext : %s\n", plaintext)
	fmt.Printf("keyID     : %d\n", keyID)
	fmt.Printf("version   : %d\n\n", version)

	ct, tweak := fpeEncrypt(plaintext, keyID, version)

	fmt.Printf("\n── result ───────────────────────────────────────────\n")
	fmt.Printf("ciphertext : %s  (digits only, format preserved, untouched)\n", ct)
	fmt.Printf("tweak      : %x\n", tweak)

	fmt.Println("\n════ DECRYPT ════════════════════════════════════════")
	fmt.Printf("ciphertext : %s\n", ct)
	fmt.Printf("tweak      : %x\n\n", tweak)

	recovered, extractedKeyID, extractedVersion := fpeDecrypt(ct, tweak)

	fmt.Printf("\n── result ───────────────────────────────────────────\n")
	fmt.Printf("recovered plaintext : %s\n", recovered)
	fmt.Printf("extracted keyID     : %d\n", extractedKeyID)
	fmt.Printf("extracted version   : %d\n", extractedVersion)
	fmt.Printf("match               : %v\n", recovered == plaintext)
}

/*
════ ENCRYPT ════════════════════════════════════════
plaintext : 1234567890123456
keyID     : 15
version   : 15

  step 1 — random tweak (before stego):
           0bb624fc3fdc743d
           tweak[0] = 00001011 (random)
           tweak[1] = 10110110 (random)

  step 2 — after stego encode:
           tweak[0] = 11111011  hi=15 (keyID=15)    lo=11 (rand)
           tweak[1] = 11110110  hi=15 (version=15)  lo=6 (rand)
           tweak[2:] = 24fc3fdc743d (pure random, untouched)

  step 3 — final tweak: fbf624fc3fdc743d

── result ───────────────────────────────────────────
ciphertext : 5203028944104274  (digits only, format preserved, untouched)
tweak      : fbf624fc3fdc743d

════ DECRYPT ════════════════════════════════════════
ciphertext : 5203028944104274
tweak      : fbf624fc3fdc743d

  raw tweak         : fbf624fc3fdc743d
  tweak[0]          : 11111011
  tweak[1]          : 11110110
  extracted keyID   : tweak[0] >> 4 = 15
  extracted version : tweak[1] >> 4 = 15

── result ───────────────────────────────────────────
recovered plaintext : 1234567890123456
extracted keyID     : 15
extracted version   : 15
match               : true

*/