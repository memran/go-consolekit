# Security

Hashing, encryption/decryption, and random value generation.

## Hashing

```go
func HashSHA256(value string) string
func HashSHA256Short(value string) string
func VerifySHA256(value, hash string) bool
```

- `HashSHA256` returns a 64-character hex string
- `HashSHA256Short` returns the first 16 characters
- `VerifySHA256` compares a value against a previously computed hash

```go
hash := console.HashSHA256("hello")
// "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

short := console.HashSHA256Short("hello")
// "2cf24dba5fb0a30e"

ok := console.VerifySHA256("hello", hash)
// true

console.VerifySHA256("wrong", hash)
// false
```

## Random Generation

```go
func RandomBytes(length int) ([]byte, error)
func RandomString(length int) (string, error)
```

- `RandomBytes` generates cryptographically secure random bytes
- `RandomString` generates a URL-safe base64 encoded string of the given length

```go
bytes, err := console.RandomBytes(32)
// []byte{...}

token, err := console.RandomString(32)
// "a3Bx9K2mQ7wE5nR1vY8cF0dG4hJ6kL2p"
```

## Encryption/Decryption

AES-256-GCM symmetric encryption.

```go
func Encrypt(plaintext, key string) (string, error)
func Decrypt(ciphertext, key string) (string, error)
```

- Key is hashed with SHA256 to produce a 32-byte AES key
- Output is base64-encoded (nonce + ciphertext)
- Uses AES-GCM with a random 12-byte nonce

```go
key := "my-secret-key"

encrypted, err := console.Encrypt("sensitive data", key)
// "base64encodedstring..."

decrypted, err := console.Decrypt(encrypted, key)
// "sensitive data"

// Wrong key returns an error
_, err = console.Decrypt(encrypted, "wrong-key")
// error: ciphertext authentication failed
```
