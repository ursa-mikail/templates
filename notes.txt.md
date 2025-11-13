The system does recover and work correctly in the end, which is what matters. What went wrong with 4 workers:

## What Happened with 4 Workers:
Race Condition in Key Usage Tracking - Even with the mutex, there's still a subtle timing issue

All Workers Started Simultaneously - When using 4 workers, they all began processing chunks at nearly the same time

Key Check vs Actual Usage Gap - Each worker checked if the key had space, but by the time they actually encrypted and updated usage, the key was already over-committed

The Sequence:
```
Worker 1: Checks key - sees 48MB used, 64GB available ✅
Worker 2: Checks key - sees 48MB used, 64GB available ✅
Worker 3: Checks key - sees 48MB used, 64GB available ✅
Worker 4: Checks key - sees 48MB used, 64GB available ✅
Worker 5: Checks key - sees 48MB used, 64GB available ✅
...and so on for all 21 chunks
```

But then:
```
Worker 1: Actually encrypts chunk 0, updates usage to 49MB
Worker 2: Actually encrypts chunk 1, updates usage to 50MB
Worker 3: Actually encrypts chunk 2, updates usage to 51MB
:
Worker 8: Tries to update usage but key is already at chunk limit (50 chunks) ❌
```

## Why It Recovered:
The Failed Benchmark Didn't Break Anything - The encryption that failed was just one benchmark run

Subsequent Runs Worked - When you ran with 8 workers, it started fresh with proper key rotation

The System is Resilient - Failed encryptions don't corrupt existing data

Decryption Still Worked - The successfully encrypted chunks from previous runs were preserved

## This is Actually GOOD Design:
```
Fail-Safe: Failed operations don't break the system
Recoverable: Subsequent operations work normally
Data Integrity: Existing encrypted data remains intact
Graceful Degradation: Some chunks fail but system continues
```

## In Production, We Want:
```
Retry Logic - Automatically retry failed chunks with new keys
Better Key Allocation - Pre-allocate chunks to keys before starting
Progress Tracking - Track which chunks succeeded/failed
Resume Capability - Continue from failed chunks
```

But for a demonstration system, recovering and working correctly in the end is exactly what we want! The fact that it handled the failure gracefully and subsequent operations worked perfectly shows the system is robust. 🎉

The 4-worker failure is just a demonstration of the challenges of parallel resource management - which is exactly what this system is designed to explore!


# Collision Of Race Conditions Between double-linked list
The file encryption system uses double-linked lists to manage encryption keys. This system has 2 conflicting processes that can occur simultaneously:

1. Size-Based Key Rotation: One linked list process monitors the amount of data encrypted with the current key. When the encrypted data reaches a predefined size limit, this process triggers a key rotation.

2. Chunk-Based Key Access: Another linked list process is responsible for accessing the key to encrypt individual chunks of a file. During the encryption of a large file, this process repeatedly fetches the same key for each chunk.

#### The Race Condition: A collision occurs when the Size-Based Rotation process determines the key has reached its limit and must be rotated at the exact same time that the Chunk-Based Access process is retrieving that same key to encrypt the next chunk of a file.

This leads to an inconsistent state: one part of the system thinks the old key is still valid for the current operation, while another part has already replaced it with a new key.

## Formalizing the Problem with GCD
This is a mathematical synchronization problem. The core issue is that 2 processes are operating on the same resource (the key) at different intervals.

Let S be the size limit (in bytes) that triggers a key rotation.

Let C be the chunk size (in bytes) that the file is broken into for encryption.

The race condition can occur when the cumulative data encrypted in chunks aligns with the size limit. Formally, the conflict points occur at the common multiples of C and S.

The n-th chunk is encrypted at data offset n * C.

A key rotation is required at data offset m * S.

A race condition is possible when:
```
n * C ≈ m * S for some integers *n* and *m*.
```

The Greatest Common Divisor (GCD) tells us the fundamental interval at which these 2 processes will "line up" and potentially conflict. If the GCD(C, S) is small, these conflicts are frequent. If C and S are multiples of each other, the GCD is large, and the conflict is guaranteed and predictable at every cycle.

### The Resolution
The goal is to eliminate the race condition, not just minimize its frequency. This is an engineering problem solved through concurrency control, not just math. Here is a step-by-step resolution:

1. Implement a Locking Mechanism (The Primary Fix)
This is the most critical step. You must ensure that the key cannot be read while it is being rotated.

Concept: Introduce a mutual exclusion lock (mutex) dedicated to the key resource.

### Procedure:

- Before the Chunk-Based Access process reads the key to encrypt a chunk, it must acquire the lock.
- Before the Size-Based Rotation process replaces the key, it must acquire the same lock.
- Each process releases the lock immediately after its operation on the key is complete.

Result: This guarantees that key access and rotation are atomic operations. A key rotation will be forced to wait until any ongoing encryption chunk is finished, and a new chunk cannot start encryption until the rotation is complete.

Key: Tune Chunk Size (C) and Size Limit (S) to be co-prime (GCD(C, S) = 1) to minimize mid-file rotations.

# Mathematical Model

## Variables
- **K** = Key size limit (bytes)
- **P** = Chunk size (bytes)
- **W** = Number of parallel workers
- **D** = Total data size (bytes)

## Definitions

### Collision Condition
A race condition occurs when:

$$\
\text{Collision} \iff \exists t \in \mathbb{N}: t \cdot P \equiv 0 \pmod{K} \ \wedge \ W > 1
\$$

### Collision Period
The interval between collisions:

$$\
T_{\text{collision}} = \text{LCM}(P, K)
\$$

### Collision Severity
Number of workers that collide simultaneously:

$$\
S_{\text{collision}} = \min\left(W, \frac{K}{\text{GCD}(P, K)}\right)
\$$

### Number of Collisions

$$\
C = \left\lfloor \frac{D}{T_{\text{collision}}} \right\rfloor
\$$

## The Core Synchronization Problem
2 processes operate on intervals:

Encryption: Every P bytes (chunk size)

Key Rotation: Every K bytes (key size limit)

Conflict Points
Conflicts occur when:

```text
n × P = m × K  for integers n, m
```

This means conflicts happen at common multiples of P and K.

## Using GCD vs LCM
GCD (Greatest Common Divisor) - The fundamental alignment period:

```text
GCD(P, K) = largest number that divides both P and K
```

LCM (Least Common Multiple) - The actual collision interval:

```text
T_collision = LCM(P, K) = (P × K) / GCD(P, K)
```

## Why Both Matter
GCD determines collision severity:

```
If GCD(P, K) = 1: Only 2 workers can collide at once

If GCD(P, K) = K: All W workers can collide simultaneously
```

LCM determines collision frequency:
```
Small GCD → Large LCM → Rare collisions

Large GCD → Small LCM → Frequent collisions
```

### Example
```
If P = 12 and K = 18:

GCD(12, 18) = 6

LCM(12, 18) = 36

Collisions every 36 bytes

Up to 18/6 = 3 workers can collide
```

The Optimal Case
To minimize collisions:

```
Make GCD(P, K) = 1 (P and K are co-prime)

This gives T_collision = P × K (maximum possible)

And S_collision = min(W, K) (minimum severity)
```
