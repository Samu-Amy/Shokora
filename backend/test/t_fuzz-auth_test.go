package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"
)

// Stress test con valori casuali per controllare panic/crash (da fare sulle route "complete" per verificare che i dati non creino problemi da qualche parte)

func FuzzRegisterUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzRegisterUserRoute -fuzztime=20s

	f.Add("Mario", "Rossi", "31-12", "mario@example.com", "Password%123!", "Password%123!")

	f.Fuzz(func(t *testing.T, firstName, lastName, birthday, email, password, passwordConf string) {

		// Make request
		req := makeRegisterUserReq(firstName, lastName, birthday, email, password, passwordConf)

		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", req)

		// - Checks -

		// No server errors
		if w.Code >= 500 {
			t.Fatalf("Server error: %v", req)
		}

		// Correct response and data
		if w.Code == 201 {
			var res APIResponse[payloads.RegisterUserRes]

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			// Check important data
			if res.Data.User.Email == "" {
				t.Fatalf("empty email on success:\nReq: %+v\nRes:%+v", req, res)
			}

			if res.Data.User.IsVerified || res.Data.User.Role != 0 {
				t.Fatal("user created as verified or with role != 0")
			}
		}

		if birthday != "" {
			_, err := authservice.ConvertBirthdayToTime(strings.TrimSpace(birthday))
			if err != nil && w.Code == 201 {
				t.Fatalf("invalid birthday accepted: %s", birthday)
			}
		}
	})
}

func FuzzLoginUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzLoginUserRoute -fuzztime=20s

	f.Add("mario@example.com", "Password%123!")

	f.Fuzz(func(t *testing.T, email, password string) {

		// Use some valid (generated with seeding) data
		if customRand.Float32() < 0.3 { // ~ 30% use valid email

			// Get a random user and his session
			i := customRand.Intn(len(authState.Users))
			user := authState.Users[i]

			email = user.Email
			password = user.PlainPassword
		}

		// Make request
		req := makeLoginUserReq(email, password)

		w := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/user", req)

		if w.Code >= 500 {
			t.Fatalf("Server error: %v", req)
		}

		if w.Code == 200 {

			var res APIResponse[payloads.LoginUserRes]

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			// Check important data
			if res.Data.VerificationId == nil {
				if res.Data.User == nil {
					t.Fatal("user not sent")
				}

				if res.Data.User.Email == "" {
					t.Fatalf("empty email on success:\nReq: %+v\nRes:%+v", req, res)
				}

				// Check cookies
				cookies := w.Result().Cookies()

				if len(cookies) <= 0 {
					t.Fatal("no cookies set")
				}

				var foundAccess bool
				var foundRefresh bool

				for _, c := range cookies {
					switch c.Name {

					case api.AccessTokenCookieName:
						foundAccess = true

						if c.Value == "" {
							t.Fatal("empty auth cookie")
						}

						if !c.HttpOnly {
							t.Error("auth cookie not HttpOnly")
						}

						if !c.Secure {
							t.Error("auth cookie not Secure")
						}

					case api.RefreshTokenCookieName:
						foundRefresh = true

						if c.Value == "" {
							t.Fatal("empty auth cookie")
						}

						if !c.HttpOnly {
							t.Error("auth cookie not HttpOnly")
						}

						if !c.Secure {
							t.Error("auth cookie not Secure")
						}
					}
				}

				if !foundAccess || !foundRefresh {
					t.Fatal("auth cookie not found")
				}

			} else {
				verificationType := getVerificationType(t, res.Data.VerificationId)

				if verificationType == auth.EmailVerification && res.Data.User == nil {
					t.Fatal("user not sent while email verification requested")

				} else if verificationType == auth.TwoFactorAuth && res.Data.User != nil {
					t.Fatal("user sent while 2FA active")
				}
			}
		}
	})
}

func FuzzLogoutUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzLogoutUserRoute -fuzztime=20s

	f.Add("Or71zHGPKEDE89eOyxiWZwB0yUlyC12Uoz9Xfzat3PM", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")

	f.Fuzz(func(t *testing.T, fuzzRefreshToken, fuzzAccessToken string) {

		// validCase := false

		var accessCookie http.Cookie
		var refreshCookie http.Cookie

		// Create JTW access token
		timeNow := time.Now().UTC()
		accessTokenExpiresAt := timeNow.Add(configs.Auth.Token.AccessTokenExp)
		refreshTokenExpiresAt := timeNow.Add(configs.Auth.Token.RefreshTokenExp)

		if customRand.Float32() < 0.3 {
			// Valid cookies

			// validCase = true

			// Get a random user and his session
			i := customRand.Intn(len(authState.Users))
			user := authState.Users[i]

			session, ok := authState.Sessions[user.Id]
			if !ok {
				t.Fatal("session not found")
			}

			claims := auth.UserClaims{
				UserId:    user.Id,
				SessionId: session.Id,
				RegisteredClaims: jwt.RegisteredClaims{
					// Subject:   strconv.FormatInt(userId, 10),
					ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
					IssuedAt:  jwt.NewNumericDate(timeNow),
					NotBefore: jwt.NewNumericDate(timeNow),
					Issuer:    "Shokora",
					Audience:  []string{"Shokora"},
				},
			}

			// Generate Access Token (and add claims)
			accessToken, err := testJwtAuthenticator.GenerateJWTToken(claims)
			if err != nil {
				t.Skip("error creating jwt token, skipping iteration")

			}

			// Create cookies
			accessCookie = api.NewSecureCookie(api.AccessTokenCookieName, accessToken, accessTokenExpiresAt)
			refreshCookie = api.NewSecureCookie(api.RefreshTokenCookieName, session.PlainToken, refreshTokenExpiresAt)

		} else {
			// Invalid cookies

			accessCookie = api.NewSecureCookie(api.AccessTokenCookieName, fuzzAccessToken, accessTokenExpiresAt)
			refreshCookie = api.NewSecureCookie(api.RefreshTokenCookieName, fuzzRefreshToken, refreshTokenExpiresAt)
		}

		// Create HTTP request
		req := httptest.NewRequest("GET", "/api/v1/auth/logout", nil)
		req.AddCookie(&accessCookie)
		req.AddCookie(&refreshCookie)

		// Recorder
		w := httptest.NewRecorder()

		// Call il router
		testRouter.ServeHTTP(w, req)

		// - Checks -

		// No server errors
		if w.Code >= 500 {
			t.Fatal("Server error")
		}

		// Check response
		// if validCase && w.Code >= 400 {
		// 	t.Fatalf("Valid tokens rejected: %d", w.Code)
		// }
	})
}

func FuzzGoogleCallbackRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzGoogleCallbackRoute -fuzztime=20s
	f.Add("Or71zHGPKEDE89eOyxiWZwB0yUlyC12Uoz9Xfzat3PM", "4_0AY0e_g4kQ")

	f.Fuzz(func(t *testing.T, state, code string) {

		// Make request
		req := payloads.OAuthGoogleCallbackReq{
			State: state,
			Code:  code,
		}

		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/google/callback", req)

		// - Checks -

		// No server errors
		if w.Code >= 500 {
			t.Fatalf("Server error: %v", req)
		}

		// Correct response and data
		if w.Code == 201 {
			// Check important data
			if state == "" || code == "" {
				t.Fatal("accepted empty value")
			}

			// - Without a valid code, it should only return not valid -

			// var res APIResponse[payloads.RegisterUserRes]

			// err := json.Unmarshal(w.Body.Bytes(), &res)
			// if err != nil {
			// 	t.Fatalf("failed to unmarshal response body: %v", err)
			// }

			//
			// if res.Data.User.Email == "" {
			// 	t.Fatalf("empty email on success:\nReq: %+v\nRes:%+v", req, res)
			// }

			// if res.Data.User.Role != 0 {
			// 	t.Fatal("user created with role != 0")
			// }

			// if !res.Data.User.IsVerified {
			// 	t.Errorf("user created as not verified") // usually it should be verified with Google
			// }
		}
	})
}
