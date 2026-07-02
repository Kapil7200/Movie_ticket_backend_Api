package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ContextKeyClaims = "claims"

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secret string
}

func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{secret: secret}
}

func (j *JWTUtil) CreateToken(userID uint, userName string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTUtil) ParseToken(tokenString string) (*JWTClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*JWTClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JWTUtil) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set(ContextKeyClaims, claims)
		ctx.Next()
	}
}

func GetClaims(ctx *gin.Context) *JWTClaims {
	value, exists := ctx.Get(ContextKeyClaims)
	if !exists {
		return nil
	}
	claims, ok := value.(*JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
