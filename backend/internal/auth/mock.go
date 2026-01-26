package auth

// import (
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type TestAuthenticator struct {
// }

// const testSecret = "test"

// var testClaims = jwt.MapClaims{
// 	"sub": int64(42),
// 	"exp": time.Now().Add(time.Hour).Unix(),
// 	// "iat": time.Now().Unix(),
// 	// "nbf": time.Now().Unix(),
// 	"iss": "test-iss",
// 	"aud": "test-aud",
// }

// func (a *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

// 	return token.SignedString([]byte(testSecret))
// }
// func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
// 	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
// 		return []byte(testSecret), nil
// 	})
// }
