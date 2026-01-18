package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"time"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWT struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWT(secretKey string, tokenDuration int64) *JWT {
	return &JWT{
		secretKey:     secretKey,
		tokenDuration: time.Duration(tokenDuration) * time.Hour,
	}
}

// GenerateToken Generate new JWT token
func (c *JWT) GenerateToken(userID int64, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.secretKey))
}

// ValidateToken Validate and parse JWT token
func (c *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperror.ErrInvalidToken
			}
			return []byte(c.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperror.ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, apperror.ErrTokenExpired
	}

	return claims, nil
}

// GetTokenDuration Get token expiration duration
func (c *JWT) GetTokenDuration() time.Duration {
	return c.tokenDuration
}
