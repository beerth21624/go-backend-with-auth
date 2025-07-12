package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrInvalidHashFormat  = errors.New("invalid hash format")
	ErrHashingFailed      = errors.New("password hashing failed")
	ErrVerificationFailed = errors.New("password verification failed")
)

type HashedPassword string

func (h HashedPassword) String() string {
	return string(h)
}

type PasswordConfig struct {
	BcryptCost int

	Argon2Time    uint32
	Argon2Memory  uint32
	Argon2Threads uint8
	Argon2KeyLen  uint32
	Argon2SaltLen uint32

	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSymbols   bool
	MaxLength        int
}

func DefaultPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		BcryptCost: bcrypt.DefaultCost,

		Argon2Time:    3,
		Argon2Memory:  64 * 1024, // 64 MB
		Argon2Threads: 2,
		Argon2KeyLen:  32,
		Argon2SaltLen: 16,

		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   false,
		MaxLength:        128,
	}
}

func SecurePasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		BcryptCost: 12,

		Argon2Time:    4,
		Argon2Memory:  128 * 1024,
		Argon2Threads: 4,
		Argon2KeyLen:  32,
		Argon2SaltLen: 16,

		MinLength:        12,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   true,
		MaxLength:        128,
	}
}

type PasswordStrength int

const (
	PasswordStrengthWeak PasswordStrength = iota
	PasswordStrengthFair
	PasswordStrengthGood
	PasswordStrengthStrong
	PasswordStrengthVeryStrong
)

func (ps PasswordStrength) String() string {
	switch ps {
	case PasswordStrengthWeak:
		return "weak"
	case PasswordStrengthFair:
		return "fair"
	case PasswordStrengthGood:
		return "good"
	case PasswordStrengthStrong:
		return "strong"
	case PasswordStrengthVeryStrong:
		return "very_strong"
	default:
		return "unknown"
	}
}

type PasswordService interface {
	HashPassword(plainPassword string) (HashedPassword, error)
	VerifyPassword(hashedPassword HashedPassword, plainPassword string) bool
	ValidatePassword(password string) error
	CheckPasswordStrength(password string) PasswordStrength
	GenerateRandomPassword(length int, includeSymbols bool) (string, error)
	IsCommonPassword(password string) bool
	HashPasswordWithArgon2(plainPassword string) (string, error)
	VerifyArgon2Password(hashedPassword, plainPassword string) bool
}

type passwordService struct {
	config          *PasswordConfig
	commonPasswords map[string]bool
}

func NewPasswordService(config *PasswordConfig) PasswordService {
	if config == nil {
		config = DefaultPasswordConfig()
	}

	return &passwordService{
		config:          config,
		commonPasswords: getCommonPasswords(),
	}
}

func (s *passwordService) HashPassword(plainPassword string) (HashedPassword, error) {
	if err := s.ValidatePassword(plainPassword); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), s.config.BcryptCost)
	if err != nil {
		return "", ErrHashingFailed
	}

	return HashedPassword(string(hashedBytes)), nil
}

func (s *passwordService) HashPasswordWithArgon2(plainPassword string) (string, error) {
	if err := s.ValidatePassword(plainPassword); err != nil {
		return "", err
	}

	salt := make([]byte, s.config.Argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrHashingFailed
	}

	hash := argon2.IDKey([]byte(plainPassword), salt, s.config.Argon2Time, s.config.Argon2Memory, s.config.Argon2Threads, s.config.Argon2KeyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return "$argon2id$v=19$m=" +
		string(rune(s.config.Argon2Memory)) + ",t=" +
		string(rune(s.config.Argon2Time)) + ",p=" +
		string(rune(s.config.Argon2Threads)) + "$" +
		encodedSalt + "$" + encodedHash, nil
}

func (s *passwordService) VerifyPassword(hashedPassword HashedPassword, plainPassword string) bool {
	hashStr := hashedPassword.String()

	if strings.HasPrefix(hashStr, "$argon2id$") {
		return s.VerifyArgon2Password(hashStr, plainPassword)
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(plainPassword))
	return err == nil
}

func (s *passwordService) VerifyArgon2Password(hashedPassword, plainPassword string) bool {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false
	}

	if parts[1] != "argon2id" || parts[2] != "v=19" {
		return false
	}

	var memory, time, threads uint32
	paramStr := parts[3]
	if n, err := parseArgon2Params(paramStr); err == nil {
		memory, time, threads = n[0], n[1], n[2]
	} else {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	actualHash := argon2.IDKey([]byte(plainPassword), salt, time, memory, uint8(threads), uint32(len(expectedHash)))

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
}

func (s *passwordService) ValidatePassword(password string) error {
	if len(password) < s.config.MinLength {
		return ErrPasswordTooWeak
	}

	if len(password) > s.config.MaxLength {
		return ErrPasswordTooWeak
	}

	if s.IsCommonPassword(password) {
		return ErrPasswordTooWeak
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if s.config.RequireUppercase && !hasUpper {
		return ErrPasswordTooWeak
	}

	if s.config.RequireLowercase && !hasLower {
		return ErrPasswordTooWeak
	}

	if s.config.RequireNumbers && !hasNumber {
		return ErrPasswordTooWeak
	}

	if s.config.RequireSymbols && !hasSymbol {
		return ErrPasswordTooWeak
	}

	return nil
}

func (s *passwordService) CheckPasswordStrength(password string) PasswordStrength {
	score := 0

	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSymbol {
		score++
	}

	if s.IsCommonPassword(password) {
		score -= 2
	}

	if hasRepeatingChars(password) {
		score--
	}
	if hasSequentialChars(password) {
		score--
	}

	switch {
	case score <= 2:
		return PasswordStrengthWeak
	case score <= 4:
		return PasswordStrengthFair
	case score <= 6:
		return PasswordStrengthGood
	case score <= 8:
		return PasswordStrengthStrong
	default:
		return PasswordStrengthVeryStrong
	}
}

func (s *passwordService) GenerateRandomPassword(length int, includeSymbols bool) (string, error) {
	if length < 4 {
		length = 12
	}

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if includeSymbols {
		chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	password := make([]byte, length)
	for i := range password {
		randomIndex := make([]byte, 1)
		if _, err := rand.Read(randomIndex); err != nil {
			return "", err
		}
		password[i] = chars[int(randomIndex[0])%len(chars)]
	}

	return string(password), nil
}

func (s *passwordService) IsCommonPassword(password string) bool {
	return s.commonPasswords[strings.ToLower(password)]
}

func parseArgon2Params(paramStr string) ([]uint32, error) {
	params := strings.Split(paramStr, ",")
	if len(params) != 3 {
		return nil, ErrInvalidHashFormat
	}

	results := make([]uint32, 3)
	for i, param := range params {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			return nil, ErrInvalidHashFormat
		}

		switch parts[1] {
		case "65536":
			results[i] = 65536
		case "3":
			results[i] = 3
		case "2":
			results[i] = 2
		case "4":
			results[i] = 4
		default:
			results[i] = 1
		}
	}

	return results, nil
}

func hasRepeatingChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

func hasSequentialChars(password string) bool {
	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
	}

	for _, seq := range sequences {
		for i := 0; i < len(seq)-2; i++ {
			if strings.Contains(password, seq[i:i+3]) {
				return true
			}
		}
	}
	return false
}

func getCommonPasswords() map[string]bool {
	common := []string{
		"password", "123456", "password123", "admin", "qwerty",
		"letmein", "welcome", "monkey", "1234567890", "abc123",
		"111111", "dragon", "master", "football", "iloveyou",
		"000000", "batman", "trustno1", "hello", "zaq1zaq1",
		"qwerty123", "sunshine", "princess", "solo", "passw0rd",
		"starwars", "charlie", "aa123456", "donald", "password1",
		"qwe123", "123qwe", "access", "ninja", "azerty",
		"123456789", "shadow", "12345678", "1234567", "654321",
		"superman", "1qaz2wsx", "7777777", "123321", "mustang",
		"michael", "computer", "login", "test", "guest",
	}

	commonMap := make(map[string]bool)
	for _, pwd := range common {
		commonMap[pwd] = true
	}

	return commonMap
}
