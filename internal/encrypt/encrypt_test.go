package encrypt_test

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/encrypt"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plain := "super-secret-value"
	pass := "my-passphrase"

	enc, err := encrypt.Encrypt(plain, pass)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.HasPrefix(enc, "enc:") {
		t.Errorf("expected enc: prefix, got %q", enc)
	}

	dec, err := encrypt.Decrypt(enc, pass)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec != plain {
		t.Errorf("expected %q, got %q", plain, dec)
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	plain, pass := "value", "pass"
	a, _ := encrypt.Encrypt(plain, pass)
	b, _ := encrypt.Encrypt(plain, pass)
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := encrypt.Encrypt("val", "")
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	enc, _ := encrypt.Encrypt("secret", "correct")
	_, err := encrypt.Decrypt(enc, "wrong")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_MissingPrefix(t *testing.T) {
	_, err := encrypt.Decrypt("plainvalue", "pass")
	if err == nil {
		t.Error("expected error for value without enc: prefix")
	}
}

func TestIsEncrypted(t *testing.T) {
	enc, _ := encrypt.Encrypt("x", "p")
	if !encrypt.IsEncrypted(enc) {
		t.Error("expected IsEncrypted true for encrypted value")
	}
	if encrypt.IsEncrypted("plain") {
		t.Error("expected IsEncrypted false for plain value")
	}
}

func TestApplyToEnv_DecryptsValues(t *testing.T) {
	pass := "testpass"
	encVal, _ := encrypt.Encrypt("db-password", pass)

	env := map[string]string{
		"DB_PASS": encVal,
		"APP_ENV": "production",
	}

	out, errs := encrypt.ApplyToEnv(env, pass)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["DB_PASS"] != "db-password" {
		t.Errorf("expected decrypted value, got %q", out["DB_PASS"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("plain value should be unchanged, got %q", out["APP_ENV"])
	}
}

func TestApplyToEnv_WrongPassphraseCollectsError(t *testing.T) {
	encVal, _ := encrypt.Encrypt("secret", "correct")
	env := map[string]string{"KEY": encVal}

	_, errs := encrypt.ApplyToEnv(env, "wrong")
	if _, ok := errs["KEY"]; !ok {
		t.Error("expected error for KEY with wrong passphrase")
	}
}
