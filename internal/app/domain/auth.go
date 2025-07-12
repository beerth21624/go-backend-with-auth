package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidJWT         = errors.New("invalid JWT token")
	ErrInvalidIPAddress   = errors.New("invalid IP address")
	ErrInvalidTokenType   = errors.New("invalid token type")
	ErrInvalidFingerprint = errors.New("invalid device fingerprint")
)

type HashedPassword string

func NewHashedPassword(plainPassword string) (HashedPassword, error) {
	if len(plainPassword) < 8 {
		return "", ErrInvalidPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return HashedPassword(hashedBytes), nil
}

func NewHashedPasswordFromHash(hash string) (HashedPassword, error) {
	if hash == "" {
		return "", ErrInvalidPassword
	}
	return HashedPassword(hash), nil
}

func (hp HashedPassword) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(plainPassword))
	return err == nil
}

func (hp HashedPassword) String() string {
	return string(hp)
}

type JWT string

var jwtRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)

func NewJWT(token string) (JWT, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ErrInvalidJWT
	}

	if !jwtRegex.MatchString(token) {
		return "", ErrInvalidJWT
	}

	return JWT(token), nil
}

func (j JWT) String() string {
	return string(j)
}

func (j JWT) IsEmpty() bool {
	return string(j) == ""
}

type IPAddress string

func NewIPAddress(ip string) (IPAddress, error) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return "", ErrInvalidIPAddress
	}

	if net.ParseIP(ip) == nil {
		return "", ErrInvalidIPAddress
	}

	return IPAddress(ip), nil
}

func (ip IPAddress) String() string {
	return string(ip)
}

func (ip IPAddress) IsPrivate() bool {
	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
	}

	for _, cidr := range privateRanges {
		_, privateNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if privateNet.Contains(parsedIP) {
			return true
		}
	}

	return false
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

func NewTokenType(tokenType string) (TokenType, error) {
	tt := TokenType(strings.ToLower(strings.TrimSpace(tokenType)))
	switch tt {
	case TokenTypeAccess, TokenTypeRefresh:
		return tt, nil
	default:
		return "", ErrInvalidTokenType
	}
}

func (tt TokenType) String() string {
	return string(tt)
}

func (tt TokenType) IsAccess() bool {
	return tt == TokenTypeAccess
}

func (tt TokenType) IsRefresh() bool {
	return tt == TokenTypeRefresh
}

type DeviceFingerprint string

func NewDeviceFingerprint(fingerprint string) (DeviceFingerprint, error) {
	fingerprint = strings.TrimSpace(fingerprint)
	if len(fingerprint) < 10 || len(fingerprint) > 500 {
		return "", ErrInvalidFingerprint
	}

	return DeviceFingerprint(fingerprint), nil
}

func GenerateDeviceFingerprint(userAgent, ip string) (DeviceFingerprint, error) {
	if userAgent == "" || ip == "" {
		return "", ErrInvalidFingerprint
	}
	combined := userAgent + "|" + ip

	encoded := base64.StdEncoding.EncodeToString([]byte(combined))

	return NewDeviceFingerprint(encoded)
}

func (df DeviceFingerprint) String() string {
	return string(df)
}

type RefreshTokenValue string

func NewRefreshTokenValue(value string) (RefreshTokenValue, error) {
	value = strings.TrimSpace(value)
	if len(value) < 32 {
		return "", ErrInvalidJWT
	}

	return RefreshTokenValue(value), nil
}

func GenerateRefreshTokenValue() (RefreshTokenValue, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	value := base64.URLEncoding.EncodeToString(bytes)
	return RefreshTokenValue(value), nil
}

func (rtv RefreshTokenValue) String() string {
	return string(rtv)
}

func (rtv RefreshTokenValue) IsEmpty() bool {
	return string(rtv) == ""
}

type TokenClaims struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func (tc *TokenClaims) IsExpired() bool {
	return time.Now().After(tc.ExpiresAt)
}

func (tc *TokenClaims) IsAccessToken() bool {
	return tc.TokenType == string(TokenTypeAccess)
}

func (tc *TokenClaims) IsRefreshToken() bool {
	return tc.TokenType == string(TokenTypeRefresh)
}
