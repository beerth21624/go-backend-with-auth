package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenNotValidYet  = errors.New("token not valid yet")
	ErrInvalidSigningKey = errors.New("invalid signing key")
	ErrInvalidKeyFormat  = errors.New("invalid key format")
)

type JWT string

func (j JWT) String() string {
	return string(j)
}

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type JWTClaims struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	TokenType   string `json:"token_type"`
	SessionID   int64  `json:"session_id,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	PrivateKey           *rsa.PrivateKey
	PublicKey            *rsa.PublicKey
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	Audience             string
}

func DefaultJWTConfig() (*JWTConfig, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &JWTConfig{
		PrivateKey:           privateKey,
		PublicKey:            &privateKey.PublicKey,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour, // 7 days
		Issuer:               "venturex-backend",
		Audience:             "venturex-app",
	}, nil
}

func LoadJWTConfigFromPEM(privateKeyPEM, publicKeyPEM string) (*JWTConfig, error) {
	privateBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privateBlock == nil {
		return nil, ErrInvalidKeyFormat
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		parsedKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, ErrInvalidKeyFormat
		}
	}

	publicBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if publicBlock == nil {
		return nil, ErrInvalidKeyFormat
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKeyFormat
	}

	return &JWTConfig{
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               "venturex-backend",
		Audience:             "venturex-app",
	}, nil
}

type JWTService interface {
	GenerateAccessToken(userID int64, username, email string, sessionID int64, fingerprint string) (JWT, time.Time, error)
	GenerateRefreshToken(userID int64, username, email string, sessionID int64) (JWT, time.Time, error)
	ValidateToken(token string) (*JWTClaims, error)
	ValidateAccessToken(token string) (*JWTClaims, error)
	ValidateRefreshToken(token string) (*JWTClaims, error)
	RefreshAccessToken(refreshToken string) (JWT, time.Time, error)
	IsTokenExpired(token string) bool
	GetTokenClaims(token string) (*JWTClaims, error)
}

type jwtService struct {
	config *JWTConfig
}

func NewJWTService(config *JWTConfig) JWTService {
	return &jwtService{config: config}
}

func (s *jwtService) GenerateAccessToken(userID int64, username, email string, sessionID int64, fingerprint string) (JWT, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.AccessTokenDuration)

	claims := &JWTClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		TokenType:   string(string(TokenTypeAccess)),
		SessionID:   sessionID,
		Fingerprint: fingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{s.config.Audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-%d", userID, sessionID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.config.PrivateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return JWT(tokenString), expiresAt, nil
}

func (s *jwtService) GenerateRefreshToken(userID int64, username, email string, sessionID int64) (JWT, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.RefreshTokenDuration)

	claims := &JWTClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		TokenType: string(string(TokenTypeRefresh)),
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{s.config.Audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("refresh-%d-%d", userID, sessionID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.config.PrivateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return JWT(tokenString), expiresAt, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.PublicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != string(string(TokenTypeAccess)) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != string(string(TokenTypeRefresh)) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *jwtService) RefreshAccessToken(refreshTokenString string) (JWT, time.Time, error) {
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	return s.GenerateAccessToken(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.SessionID,
		claims.Fingerprint,
	)
}

func (s *jwtService) IsTokenExpired(tokenString string) bool {
	_, err := s.ValidateToken(tokenString)
	return errors.Is(err, ErrTokenExpired)
}

func (s *jwtService) GetTokenClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.PublicKey, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func PrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

func PublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}
