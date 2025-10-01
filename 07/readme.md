# gRPC File Signing Service

A secure gRPC service that demonstrates TLS encryption and file signing with RSA signatures.

## 🎯 Features

- **TLS Encryption**: Secure gRPC communication with mutual TLS
- **File Signing**: Sign files using RSA-SHA256
- **Signature Verification**: Client-side verification of signatures
- **Simple Makefile**: Easy-to-use commands

## 📋 Usage
```
% make clean
🧹 Cleaning...
✅ Clean completed!

% make all  
🚀 Initializing project...
✅ Project initialized!
📦 Installing protobuf tools...
✅ Tools installed!
🔐 Generating TLS certificates...
=== Cleaning up old certificates ===

=== Step 1: Generating CA certificate ===
✓ CA certificate created

=== Step 2: Generating server certificate with SANs ===
Certificate request self-signature ok
subject=C = US, ST = California, L = San Francisco, O = Test Server, CN = localhost
✓ Server certificate created

=== Step 3: Generating client certificate with SANs ===
Certificate request self-signature ok
subject=C = US, ST = California, L = San Francisco, O = Test Client, CN = client
✓ Client certificate created

=== Cleaning up temporary files ===

=== ✅ All certificates generated successfully! ===

Final files:
-rw-r--r--  1 chanfamily  staff   2.0K Sep 30 19:17 ca.crt
-rw-------  1 chanfamily  staff   3.2K Sep 30 19:17 ca.key
-rw-r--r--  1 chanfamily  staff   2.0K Sep 30 19:17 client.crt
-rw-------  1 chanfamily  staff   3.2K Sep 30 19:17 client.key
-rw-r--r--  1 chanfamily  staff   2.0K Sep 30 19:17 server.crt
-rw-------  1 chanfamily  staff   3.2K Sep 30 19:17 server.key

✅ Certificates generated in certs/
⚙️  Generating protobuf code...
✅ Protobuf code generated!
🔨 Building server...
✅ Server built: bin/server
🔨 Building client...
✅ Client built: bin/client

✅ Complete! Project is ready to use.
Run 'make test' to verify everything works.
(base) chanfamily@Chans-MacBook-Air 07 % make test
⚙️  Generating protobuf code...
✅ Protobuf code generated!
🔨 Building server...
✅ Server built: bin/server
🔨 Building client...
✅ Client built: bin/client
🧪 Running test...
2025/09/30 19:17:39 Server listening on :50051 with TLS
2025/09/30 19:17:41 Received hello from: gRPC Client
2025/09/30 19:17:41 Greeting: Hello gRPC Client
2025/09/30 19:17:41 Sending file for signing: test.txt (32 bytes)
2025/09/30 19:17:41 Received file to sign: test.txt (32 bytes)
2025/09/30 19:17:41 File signed successfully: test.txt
2025/09/30 19:17:41 File signed successfully!
2025/09/30 19:17:41 Algorithm: SHA256-RSA-PKCS1v15
2025/09/30 19:17:41 Signature length: 256 bytes
2025/09/30 19:17:41 Public key length: 459 bytes
2025/09/30 19:17:41 Signature saved to: test.txt.sig
2025/09/30 19:17:41 Public key saved to: test.txt.pub
2025/09/30 19:17:41 Signature verified successfully!
✅ Test completed!

```

```
% make test
⚙️  Generating protobuf code...
✅ Protobuf code generated!
🔨 Building server...
✅ Server built: bin/server
🔨 Building client...
✅ Client built: bin/client
🧪 Running test...
2025/09/30 19:19:20 Server listening on :50051 with TLS
2025/09/30 19:19:22 Received hello from: gRPC Client
2025/09/30 19:19:22 Greeting: Hello gRPC Client
2025/09/30 19:19:22 Sending file for signing: test.txt (32 bytes)
2025/09/30 19:19:22 Received file to sign: test.txt (32 bytes)
2025/09/30 19:19:22 File signed successfully: test.txt
2025/09/30 19:19:22 File signed successfully!
2025/09/30 19:19:22 Algorithm: SHA256-RSA-PKCS1v15
2025/09/30 19:19:22 Signature length: 256 bytes
2025/09/30 19:19:22 Public key length: 459 bytes
2025/09/30 19:19:22 Signature saved to: test.txt.sig
2025/09/30 19:19:22 Public key saved to: test.txt.pub
2025/09/30 19:19:22 Signature verified successfully!
✅ Test completed!

```

This will:
1. Build server and client
2. Generate a test file
3. Start server in background
4. Sign the test file
5. Verify the signature
6. Clean up

