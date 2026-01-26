package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// - Context -

// Keys
type contextKey uint8

const (
	userCtx contextKey = iota
)

// Functions
func getUserFromContext(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	return user, ok
}

// - Params -

// Constants
const userIdParam = "userId"
const productIdParam = "productId"

// Methods
func (app *App) getIdFromParam(r *http.Request, idParamName string) (int64, error) {
	idParam := chi.URLParam(r, idParamName)

	resourceId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return -1, err
	}

	return resourceId, nil
}

// - Auth verification (verify email, reset password, 2FA) -

// Enum
type tokenType uint8

const (
	tokenEmailVerification tokenType = iota
	tokenPasswordReset
	tokenTwoFactorAuth
)

// Functions
func generateVerificationTokens() (string, []byte) {
	plainToken, err := generateplainVerificationToken()
	if err != nil {
		return "", nil // TODO: sistema
	}

	// hash token

	return plainToken, nil
}

// func generateOTP() string {
// 	buffer := make([]byte, 32) // 32 byte
// 	if _, err := rand.Read(buffer); err != nil {
// 		return ""
// 	}
// 	return hex.EncodeToString(buffer)
// }

func generateplainVerificationToken() (string, error) {
	buffer := make([]byte, 32)                   // 32 byte
	if _, err := rand.Read(buffer); err != nil { // TODO: sistema
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// TODO: controlla
func generatePlainOTP(digits int) (string, error) {
	if digits < 4 || digits > 8 {
		return "", errors.New("invalid otp length")
	}

	max := int(math.Pow10(digits))

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", digits, n.Int64()), nil
}

func hashToken(plainToken string) ([]byte, error) {
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	return hashedToken, nil
}

// func hashToken(token string, secret []byte) []byte {
// 	mac := hmac.New(sha256.New, secret)
// 	mac.Write([]byte(token))
// 	return mac.Sum(nil)
// }

// CONFRONTO (?)
// if !hmac.Equal(storedHash, hashToken(token, secret)) {
// 	return errors.New("invalid token")
// }

func hashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, err
}

// func (app *App) setAuthCookie(w http.ResponseWriter, token string) {
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "auth_token",
// 		Value:    token,
// 		Path:     "/",
// 		MaxAge:   int(app.config.Auth.Token.Exp.Seconds()),
// 		HttpOnly: true,
// 		Secure:   app.config.Env == "production", // true in prod
// 		SameSite: http.SameSiteStrictMode,
// 	})
// }
