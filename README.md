# FPE Tweak Steganography

Hides `keyID` and `version` inside a random tweak used for FF1 format-preserving encryption. The tweak carries metadata invisibly — no separate storage, no plaintext fields, no envelope.

---

## How the encoding works

The scheme mirrors Ubiq's `encodeKeyNumber` / `decodeKeyNumber` pattern, applied to raw bytes instead of alphabet characters.

Each byte (0–255) is split into two nibbles:

```
 7   6   5   4   3   2   1   0
└───────────────┘└───────────────┘
   high nibble       low nibble
   (hidden data)     (random noise)
```

- **High nibble** (`bits 7–4`) — carries the hidden value (`keyID` or `version`)
- **Low nibble** (`bits 3–0`) — retains the original random bits from the pre-stego tweak

`sft = 4` is the shift constant. The encode and decode functions are:

```go
encodeByte(r, n) → low = r & 0x0f;  return low + (n << 4)
decodeByte(b)    → return b >> 4
```

This directly matches the Ubiq pattern:
```
encode: idx = posOf(inp[0]);  idx += n << sft;  inp[0] = valAt(idx)
decode: c   = posOf(inp[0]);  n    = c >> sft;   inp[0] = valAt(c - (n<<sft))
```
For raw bytes, `posOf` and `valAt` are identity — the byte value is its own index.

---

## Tweak layout

The tweak is **8 bytes**, generated fully at random first, then bytes `[0]` and `[1]` are overwritten by the steganography encoder:

```
byte index:  [0]          [1]          [2]    [3]    [4]    [5]    [6]    [7]
             ┌──────────┐ ┌──────────┐ ┌────────────────────────────────────┐
             │hhhh rrrr │ │hhhh rrrr │ │         random nonce               │
             └──────────┘ └──────────┘ └────────────────────────────────────┘
              ↑    ↑        ↑    ↑
              │    └─ 4 random bits (from pre-stego random byte, low nibble)
              └─ keyID      └─ version
                 (4 bits,      (4 bits,
                  high nibble)  high nibble)
```

### Step by step

```
step 1 — generate 8 random bytes:
         e.g. b8c638c5c669f80b

step 2 — encode keyID into tweak[0]:
         tweak[0] was: 1011 1000  (random)
         keyID = 7:    0111
         encoded:      0111 1000  → high nibble = 7, low nibble = 8 (rand)

         encode version into tweak[1]:
         tweak[1] was: 1100 0110  (random)
         version = 3:  0011
         encoded:      0011 0110  → high nibble = 3, low nibble = 6 (rand)

step 3 — tweak[2:8] untouched:
         38c5c669f80b  (pure random)

final tweak: 7836 38c5c669f80b
             ↑↑
             tweak[0]: keyID=7 in high nibble, rand=8 in low nibble
               tweak[1]: version=3 in high nibble, rand=6 in low nibble
```

### Extraction on decrypt

```go
keyID   = tweak[0] >> 4   // strip low nibble, read high nibble
version = tweak[1] >> 4
```

No mutation of the tweak is needed — the same tweak bytes go directly into
`ff1.NewCipher` for both encrypt and decrypt.

---

## Value ranges

| Field     | Tweak byte | Bits used | Range  | Constraint                      |
|-----------|------------|-----------|--------|---------------------------------|
| `keyID`   | `[0]`      | high 4    | 0–15   | `(15 << 4) + 15 = 255 ≤ 255` ✓ |
| `version` | `[1]`      | high 4    | 0–15   | same ✓                          |
| nonce     | `[2:8]`    | all 8     | —      | 6 bytes pure random             |

Maximum values:
- `keyID`:   **15** (4 bits)
- `version`: **15** (4 bits)

The hard limit is `(n << sft) + 15 ≤ 255`, which gives `n ≤ 15` for `sft = 4`.

---

## Extending the range

To support larger values without changing `sft` (preserving 4 bits of noise per byte),
spread the data across more bytes:

| Goal                | Layout                                                         | keyID range |
|---------------------|----------------------------------------------------------------|-------------|
| Current (this code) | `[0]`=keyID, `[1]`=version                                     | 0–15        |
| keyID 0–255         | `[0]`=keyID high nibble, `[1]`=keyID low nibble, `[2]`=version | 0–255       |
| keyID 0–65535       | `[0:4]` carry keyID nibble-by-nibble, `[4]`=version            | 0–65535     |

Each additional tweak byte buys 4 more bits of keyID range while keeping
4 bits of random noise per byte throughout.
