#!/bin/bash

# Exit on any error
set -e

# Generate TLS certificates for gRPC mutual TLS authentication

echo "=== Cleaning up old certificates ==="
rm -f *.crt *.key *.csr *.srl *.cnf

echo ""
echo "=== Step 1: Generating CA certificate ==="
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days 365 -key ca.key -out ca.crt \
    -subj "/C=US/ST=California/L=San Francisco/O=Test CA/CN=Test CA"
echo "✓ CA certificate created"

echo ""
echo "=== Step 2: Generating server certificate with SANs ==="
openssl genrsa -out server.key 4096

# Create server config with SANs
cat > server.cnf <<EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = US
ST = California
L = San Francisco
O = Test Server
CN = localhost

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

openssl req -new -key server.key -out server.csr -config server.cnf
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out server.crt -extensions req_ext -extfile server.cnf
echo "✓ Server certificate created"

echo ""
echo "=== Step 3: Generating client certificate with SANs ==="
openssl genrsa -out client.key 4096

# Create client config with SANs
cat > client.cnf <<EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = US
ST = California
L = San Francisco
O = Test Client
CN = client

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = client
DNS.2 = localhost
EOF

openssl req -new -key client.key -out client.csr -config client.cnf
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out client.crt -extensions req_ext -extfile client.cnf
echo "✓ Client certificate created"

echo ""
echo "=== Cleaning up temporary files ==="
rm -f *.csr *.cnf

echo ""
echo "=== ✅ All certificates generated successfully! ==="
echo ""
echo "Final files:"
ls -lh *.crt *.key
echo ""