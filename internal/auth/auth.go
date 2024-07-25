package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type CustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func CreateToken(data int, expiresAt time.Duration, key string) (string, error) {
	claims := CustomClaims{
		data,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Unix(int64(expiresAt), 0)),
			Issuer:    "chirpy",
			Subject:   strconv.Itoa(data),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString, secretKey string) (*jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(*CustomClaims); ok && tkn.Valid {
		fmt.Printf("Token is valid. User ID: %v\n", claims.Subject)
		// Further validation can be done here, such as issuer check
		return tkn, nil
	} else {
		// log.Fatal("unknown claims type, cannot proceed")
		return nil, fmt.Errorf("invalid token")
	}
}

func GetToken(r *http.Request) (string, error) {
	signedToken := r.Header.Get("Authorization")

	if signedToken == "" {
		return "", fmt.Errorf("no token provided")
	}

	newToken := strings.Split(signedToken, " ")

	if newToken[0] != "Bearer" || len(newToken) != 2 {
		return "", fmt.Errorf("invalid token")
	}

	return newToken[1], nil
}

func CreateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	refreshToken := hex.EncodeToString(randomBytes)

	return refreshToken, nil
}
