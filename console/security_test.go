package console

import (
	"testing"
)

func TestHashSHA256(t *testing.T) {
	h := HashSHA256("hello")
	if len(h) != 64 {
		t.Fatalf("expected 64 hex chars, got %d: %s", len(h), h)
	}
	if h == HashSHA256("world") {
		t.Fatal("different inputs should produce different hashes")
	}
}

func TestVerifySHA256(t *testing.T) {
	h := HashSHA256("hello")
	if !VerifySHA256("hello", h) {
		t.Fatal("verify should match")
	}
	if VerifySHA256("wrong", h) {
		t.Fatal("verify should not match wrong input")
	}
}

func TestHashSHA256Short(t *testing.T) {
	h := HashSHA256Short("hello")
	if len(h) != 16 {
		t.Fatalf("expected 16, got %d", len(h))
	}
}

func TestRandomBytes(t *testing.T) {
	b1, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("random bytes failed: %v", err)
	}
	b2, _ := RandomBytes(16)
	if len(b1) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(b1))
	}
	if string(b1) == string(b2) {
		t.Fatal("expected different random bytes")
	}
}

func TestRandomString(t *testing.T) {
	s1, err := RandomString(16)
	if err != nil {
		t.Fatalf("random string failed: %v", err)
	}
	s2, _ := RandomString(16)
	if len(s1) != 16 {
		t.Fatalf("expected 16, got %d", len(s1))
	}
	if s1 == s2 {
		t.Fatal("expected different random strings")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := "my-secret-key-123"
	original := "hello world"
	encrypted, err := Encrypt(original, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if decrypted != original {
		t.Fatalf("expected '%s', got '%s'", original, decrypted)
	}
}

func TestEncryptWrongKey(t *testing.T) {
	encrypted, err := Encrypt("secret", "key1")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	_, err = Decrypt(encrypted, "key2")
	if err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
}

func TestEncryptEmpty(t *testing.T) {
	encrypted, err := Encrypt("", "key")
	if err != nil {
		t.Fatalf("encrypt empty failed: %v", err)
	}
	decrypted, err := Decrypt(encrypted, "key")
	if err != nil {
		t.Fatalf("decrypt empty failed: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("expected empty, got '%s'", decrypted)
	}
}

func TestEncryptDecryptLongText(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog. This is a longer text to test encryption."
	encrypted, err := Encrypt(text, "long-key-here-12345")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	decrypted, err := Decrypt(encrypted, "long-key-here-12345")
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if decrypted != text {
		t.Fatal("decrypted text does not match original")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	_, err := Decrypt("not-base64!!!", "key")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}
