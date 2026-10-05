package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashSessionToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	return tokenHash
}

func GenerateAccessCode() string {
	code := strings.ToUpper(uuid.NewString()[:13])
	code = code[0:4] + "-" + code[4:]

	return code
}

func NormalizeAccessCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func HashRegistrationCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])
	return codeHash
}

func ValidateTokenLength(token string) bool {
	if len(token) != 64 {
		return false
	}

	_, err := hex.DecodeString(token)
	return err == nil
}

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024 // 64 MiB, measured in KiB
	argonThreads uint8  = 2
	argonKeyLen  uint32 = 32
	saltLength          = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	encodedPassword := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonTime,
		argonThreads,
		encodedSalt,
		encodedHash,
	)

	return encodedPassword, nil
}

func PasswordValidateLength(password string) bool {
	if len(password) < 8 || len(password) > 255 {
		return false
	}
	return true
}

func RunDummyArgon2Hash(password string) (bool, error) {
	passwordHash := "$argon2id$v=19$m=65536,t=3,p=2$Xmtvvc9bm1z0Fg1XAwMDnw$oDE9hEgCcT2bnwOPXHXG9T1EcxBxSq+fGFdG6+VpOgs"

	match, err := argon2id.ComparePasswordAndHash(password, passwordHash)
	return match, err
}
