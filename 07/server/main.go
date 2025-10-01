package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "hellogrpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
	privateKey *rsa.PrivateKey
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received hello from: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SignFile(ctx context.Context, in *pb.SignFileRequest) (*pb.SignFileResponse, error) {
	log.Printf("Received file to sign: %s (%d bytes)", in.GetFileName(), len(in.GetFileContent()))
	
	// Hash the file content
	hashed := sha256.Sum256(in.GetFileContent())
	
	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	
	// Get public key in PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&s.privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	
	log.Printf("File signed successfully: %s", in.GetFileName())
	
	return &pb.SignFileResponse{
		Signature: signature,
		PublicKey: publicKeyPEM,
		Algorithm: "SHA256-RSA-PKCS1v15",
	}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return nil, err
	}
	
	// Load CA certificate
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, err
	}
	
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate")
	}
	
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}), nil
}

func generateKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func main() {
	// Generate server key pair for signing
	privateKey, err := generateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Set up TLS credentials
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	
	// Create gRPC server with TLS
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	
	// Register service
	pb.RegisterGreeterServer(grpcServer, &server{
		privateKey: privateKey,
	})
	
	// Start server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	log.Printf("Server listening on %s with TLS", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}