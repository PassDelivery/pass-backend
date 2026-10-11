package core_jwt

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uuid.UUID
	RoleUser domain.Role
}

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(ttl time.Duration) *TokenService {
	secret := getSecret()
	return &TokenService{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *TokenService) Generate(userId uuid.UUID, role domain.Role) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   userId,
		RoleUser: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signed string by secret: %w", err)
	}
	return signedToken, nil
}

func (s *TokenService) VerifyToken(token string) (*Claims, error) {
	claims := &Claims{}

	keyFunc := jwt.Keyfunc(func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("verify token: %w", core_errors.ErrInvalidArgument)
		}
		return []byte(s.secret), nil
	})

	jwtToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("get data from token: %w", core_errors.ErrInvalidArgument)
	}

	return claims, nil
}

func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("secret is not set")
	}

	return []byte(secret)
}

func FromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value("claims").(*Claims)
	return c, ok
}

func CheckPassword(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))

	return err == nil
}

func GetIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims, ok := FromContext(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("jwt from context", core_errors.ErrUnauthtorize)
	}

	return claims.UserID, nil
}
