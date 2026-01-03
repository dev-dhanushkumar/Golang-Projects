package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"not null" json:"-"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	PhoneNumber     string     `json:"phone_number"`
	ProfileImageURL string     `json:"profile_image_url"`
	DefaultCurrency string     `gorm:"default:'USD'" json:"default_currency"`
	EmailVerified   bool       `gorm:"default:false" json:"email_verified"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// BeforeCreate hook - automatically called before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	// Additional validations can be added here
	return nil
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *User) EncryptPhone(plainText string, str_key string) error {
	key := mustDecodeBase64Key(str_key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	number := base64.StdEncoding.EncodeToString(ciphertext)

	u.PhoneNumber = string(number)
	return nil
}

func mustDecodeBase64Key(encodedKey string) []byte {
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		log.Fatalf("Fatal error decoding Base64 key: %v", err)
	}
	if len(key) != 32 {
		log.Fatalf("Invalid key length: expected 32 bytes (AES-256), got %d bytes", len(key))
	}
	return key
}
