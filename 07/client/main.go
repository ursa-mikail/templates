package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "hellogrpc"
)

const (
	address = "localhost:50051"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load CA certificate
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, err
	}
	
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate")
	}
	
	// Load client certificate
	clientCert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		return nil, err
	}
	
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}), nil
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return io.ReadAll(file)
}

func savePublicKey(publicKey []byte, filename string) error {
	return os.WriteFile(filename, publicKey, 0644)
}

func saveSignature(signature []byte, filename string) error {
	return os.WriteFile(filename, signature, 0644)
}

func verifySignature(fileContent, signature, publicKeyPEM []byte) error {
	// Parse public key
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing public key")
	}
	
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("expected RSA public key")
	}
	
	// Hash the file content
	hashed := sha256.Sum256(fileContent)
	
	// Verify signature
	return rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature)
}

func main() {
	// Load TLS credentials
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	
	// Set up connection with TLS
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	
	client := pb.NewGreeterClient(conn)
	
	// Test Hello call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Say Hello
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "gRPC Client"})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	
	// Sign file if provided
	if len(os.Args) > 1 {
		filename := os.Args[1]
		
		// Read file to sign
		fileContent, err := readFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		
		log.Printf("Sending file for signing: %s (%d bytes)", filename, len(fileContent))
		
		// Send file to server for signing
		signCtx, signCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer signCancel()
		
		signResponse, err := client.SignFile(signCtx, &pb.SignFileRequest{
			FileContent: fileContent,
			FileName:    filename,
		})
		if err != nil {
			log.Fatalf("Could not sign file: %v", err)
		}
		
		log.Printf("File signed successfully!")
		log.Printf("Algorithm: %s", signResponse.GetAlgorithm())
		log.Printf("Signature length: %d bytes", len(signResponse.GetSignature()))
		log.Printf("Public key length: %d bytes", len(signResponse.GetPublicKey()))
		
		// Save signature and public key
		signatureFile := filename + ".sig"
		publicKeyFile := filename + ".pub"
		
		if err := saveSignature(signResponse.GetSignature(), signatureFile); err != nil {
			log.Fatalf("Failed to save signature: %v", err)
		}
		
		if err := savePublicKey(signResponse.GetPublicKey(), publicKeyFile); err != nil {
			log.Fatalf("Failed to save public key: %v", err)
		}
		
		log.Printf("Signature saved to: %s", signatureFile)
		log.Printf("Public key saved to: %s", publicKeyFile)
		
		// Verify signature locally
		err = verifySignature(fileContent, signResponse.GetSignature(), signResponse.GetPublicKey())
		if err != nil {
			log.Printf("Warning: Signature verification failed: %v", err)
		} else {
			log.Printf("Signature verified successfully!")
		}
	} else {
		log.Println("No file provided for signing. Usage: client <filename>")
	}
}