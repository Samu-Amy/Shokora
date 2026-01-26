package api

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
	tokenEmailVerification tokenType = 0
	tokenPasswordReset     tokenType = 1
	tokenTwoFactorAuth     tokenType = 2
)

// Generation Methods
func (app *App) generateVerificationToken() (string, error) { // TODO: nell'handler gestire il retry nel caso non dovesse essere unico
	buffer := make([]byte, app.config.Auth.MagicLink.ByteSize)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

// TODO: per la verifica dell'otp si usa anche lo user_id nella richiesta (l'otp potrebbe non essere unico nel db)
func (app *App) generateOTP() (string, error) {
	length := app.config.Auth.OTP.Length

	// Max for length = 6 -> 1000000 (values in range[000000, 999999])
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) // Create a new *big.Int as 10^length (big.NewInt(10) ^ big.NewInt(int64(length)))

	otp, err := rand.Int(rand.Reader, max) // max is a *big.Int
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, otp), nil // Format with [length] numbers/zeros
}

// Hash-related Functions/Methods
func hashToken(plainToken string) []byte {
	hash := sha256.Sum256([]byte(plainToken)) // TODO: aggiungere pepper (secret)?
	return hash[:]                            // From [32]byte to []byte
}

func verifyToken(plainToken string, hashedToken []byte) bool {
	hash := hashToken(plainToken)
	return bytes.Equal(hash, hashedToken)
}

func (app *App) hashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), app.config.Auth.HashingCost)
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
