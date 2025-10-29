
## 🔧 Key Operations
Create Key
```go
key, err := manager.DataLimitList.CreateKey(AlgoAESGCM)
```

Force New Key
```go
newKey, err := manager.DataLimitList.ForceNewKey(AlgoChaCha20)
```

Rotate Key
```go
newKey, err := manager.DataLimitList.RotateKey(oldKeyID, AlgoAESGCM)
```

Revoke Key
```go
err := manager.DataLimitList.RevokeKey(keyID)
```

The value $$\ \max(count_{access}(K)) - \min(count_{access}(K)) \$$ is kept small to prevent hotspots (preventing any single key from being a leakage or wear hotspot).

Check Available Space
```go
available, err := manager.DataLimitList.GetAvailableBytes(keyID)
```

🚨 Error Handling
The system provides comprehensive error handling for:
- Key exhaustion (bytes/chunks)
- HMAC verification failures
- File I/O errors
- Concurrent access violations
- Cryptographic operation failures

📈 Performance
The system automatically benchmarks different worker counts:

```text
Workers  Time        Throughput  Chunks/Sec    File Size
1        125ms       16.02 MB/s  3.0           2.00 MB
2        78ms        25.64 MB/s  4.8           2.00 MB  
4        45ms        44.44 MB/s  8.3           2.00 MB
```

🔒 Security Features
- Key Isolation: Each chunk can use different encryption keys
- Usage Limits: Keys automatically rotate when limits reached
- Integrity Checks: HMAC verification on all encrypted chunks
- Secure Randomness: Cryptographically secure random number generation
- Key Revocation: Immediate key invalidation support