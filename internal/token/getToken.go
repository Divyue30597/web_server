package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func CreateToken(data int, expiresAt int64, key string) (string, error) {
	claims := CustomClaims{
		data,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0).UTC()),
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
