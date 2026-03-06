package auth

import (
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	audience string
	issuer   string
}

// TODO: gestisci anche Refresh Tokens (?) - aggiungi pepper per l'hashing dei Refresh Token

// - Constructor -

func NewJWTAuthenticator(secret, audience, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, audience, issuer}
}

// - Methods -

func (a *JWTAuthenticator) GenerateJWTToken(claims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate a JWT Token
func (a *JWTAuthenticator) ValidateJWTToken(token string) (*UserClaims, error) {
	if token == "" {
		return nil, interrors.IErrInvalid
	}

	claims := &UserClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		// Check signing method (algoritm)
		// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		// }

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	return claims, nil
}
