package conf

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ekstrak token dari header
		token := extractToken(c)
		if token == "" {
			unauthorized(c, "No token provided")
			return
		}

		// Parse dan validasi token
		claims, err := parseToken(token)
		if err != nil {
			unauthorized(c, err.Error())
			return
		}

		// Validasi tambahan pada klaim
		if err := validateTokenClaims(claims); err != nil {
			unauthorized(c, err.Error())
			return
		}

		// Set klaim di context untuk digunakan di handler selanjutnya
		c.Set("user_id", claims.UserId)
		c.Set("claims", claims)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func parseToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_KEY), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func validateTokenClaims(claims *JWTClaim) error {
	// Validasi waktu kedaluwarsa
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("token has expired")
	}

	// Validasi issuer (opsional)
	if claims.Issuer != "koriebruh" {
		return errors.New("invalid token issuer")
	}

	return nil
}

func unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":   http.StatusUnauthorized,
		"status": "Unauthorized",
		"error":  message,
	})
	c.Abort() // Hentikan eksekusi middleware atau handler selanjutnya
}
