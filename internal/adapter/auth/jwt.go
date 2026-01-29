package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Ensure secret exists on startup (safety check)
func init() {
	if len(jwtSecret) == 0 {
		fmt.Println("⚠️  WARNING: JWT_SECRET not set. Using insecure default for dev.")
		jwtSecret = []byte("dev-secret-do-not-use-in-prod")
	}
}

type Claims struct {
	PlayerID uuid.UUID `json:"pid"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for a player
func GenerateToken(playerID uuid.UUID) (string, error) {
	claims := Claims{
		PlayerID: playerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24h Session
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "galaxies-burn-rate",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken parses the token and extracts the Player UUID
func ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.PlayerID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}
