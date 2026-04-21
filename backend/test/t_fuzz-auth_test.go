package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
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

		// Random ID
		workerID := rand.Int63()

		seedUsersFuzz(t, db, workerID)

		// Use some valid (generated with seeding) data
		if customRand.Float32() < 0.3 { // ~ 30% use valid email
			i := customRand.Intn(min(seedUserNum, len(validEmails), len(validPasswords)))

			emailParts := strings.Split(validEmails[i], "@")
			email = fmt.Sprintf("%s-%d-%d@%s", emailParts[0], workerID, time.Now().UnixNano(), emailParts[1])
			t.Log("Valid Email tested")

			if customRand.Float32() < 0.5 { // ~ 50% use valid password for the chosen email
				password = validPasswords[i]
				t.Log("Valid Password tested")
			}
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

	f.Fuzz(func(t *testing.T, data []byte) {

		// - Generate random number for case -
		caseNum := customRand.Intn(5)

		var (
			needLogin = true
			authMode  int // 0 valid, 1 invalid, 2 none, 3 patial, 4 partial
			expected  int
		)

		switch caseNum {
		case 1: // No Login
			needLogin = false
			expected = http.StatusUnauthorized

		case 2: // Login and valid cookies
			authMode = 0
			expected = http.StatusNoContent

		case 3: // Login and not valid cookies
			authMode = 1
			expected = http.StatusUnauthorized

		case 4: // Login and no cookies
			authMode = 2
			expected = http.StatusUnauthorized

		case 5: // Login and only access token
			authMode = 3
			expected = http.StatusUnauthorized

		case 6: // Login and only refresh token
			authMode = 4
			expected = http.StatusUnauthorized

		}

		// - Login User -

		var cookies []*http.Cookie

		if needLogin {
			// Random ID
			workerID := rand.Int63()

			seedUsersFuzz(t, db, workerID)

			// Generate random data and create request
			emailParts := strings.Split(randomFrom(validEmails), "@")
			email := fmt.Sprintf("%s-%d@%s", emailParts[0], workerID, emailParts[1])

			loginReq := makeLoginUserReq(email, randomFrom(validPasswords))

			loginW := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/user", loginReq)
			cookies = loginW.Result().Cookies()

			if loginW.Code != 200 {
				t.Skip("login failed, skip iteration")
			}
		}

		// - Logout user -

		// Make request
		req := httptest.NewRequest("GET", "/api/v1/auth/logout", nil)

		switch authMode {

		case 0: // valid
			for _, c := range cookies {
				req.AddCookie(c)
			}

		case 1: // invalid
			for _, c := range cookies {
				c.Value = "garbage"
				req.AddCookie(c)
			}

		case 3: // partial
			if len(cookies) > 0 {
				req.AddCookie(cookies[0])
			}

		case 4: // partial
			if len(cookies) >= 1 {
				req.AddCookie(cookies[1])
			}
		}

		// Recorder
		w := httptest.NewRecorder()

		// Chiama il router
		testRouter.ServeHTTP(w, req)

		// - Checks -

		// No server errors
		if w.Code >= 500 {
			t.Fatal("Server error")
		}

		// Case no user/cookies
		// TODO: finisci
		// if w.Code != http.StatusUnauthorized {
		// 	t.Errorf("Expected unauthorized, got %d", w.Code)
		// }
	})
}

func FuzzGoogleLoginRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzGoogleLoginRoute -fuzztime=20s

	f.Fuzz(func(t *testing.T, req []byte) {

		// Make request
		w := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/google", nil)

		// - Checks -

		// No server errors
		if w.Code >= 500 {
			t.Fatalf("Server error")
		}

		if w.Code >= 400 {
			t.Errorf("Error with code: %v", w.Code)
		}

		// Correct response and data
		if w.Code == 201 {
			var res APIResponse[payloads.OAuthGoogleLoginRes]

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			// Check important data
			if res.Data.Url == "" {
				t.Fatalf("empty email on success:\nRes:%+v", res)
			}

			parsedURL, err := url.ParseRequestURI(res.Data.Url)
			if err != nil {
				t.Fatalf("invalid URL returned: %s\nErr: %v", res.Data.Url, err)
			}

			if parsedURL.Scheme != "https" {
				t.Fatalf("non-https URL returned: %s", res.Data.Url)
			}

			if parsedURL.Host == "" {
				t.Fatalf("URL without host: %s", res.Data.Url)
			}
		}
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
