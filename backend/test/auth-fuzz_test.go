package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
)

// Stress test con valori casuali per controllare panic/crash (da fare sulle route "complete" per verificare che i dati non creino problemi da qualche parte)

func FuzzRegisterUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzRegisterUserRoute -fuzztime=20s

	clearTestDB(db) // Already done by main (but called also here for security)

	f.Add("Mario", "Rossi", "31-12", "mario@example.com", "Password%123!", "Password%123!")

	f.Fuzz(func(t *testing.T, firstName, lastName, birthday, email, password, passwordConf string) {
		req := makeRegisterUserReq(firstName, lastName, birthday, email, password, passwordConf)

		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", req)

		// - Check errors -

		if w.Code >= 500 {
			t.Fatalf("Server error: %v", req)
		}

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

// TODO: fixxa questo test (non runna con fuzz)
func FuzzLoginUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzLoginUserRoute -fuzztime=20s

	// customRand := rand.New(rand.NewSource(randSeed))

	// configs := appconfig.NewTestConfig()
	// dataValidator := payloads.NewValidator()

	// mailer := appconfig.GetMailerFromConfig(configs)
	// jwtAuthenticator := appconfig.GetJWTAuthenticatorFromConfig(configs)
	// tokenAuthenricator := appconfig.GetTokenAuthenticatorFromConfig(configs)

	// db, err := appconfig.GetDbFromConfig(configs)
	// if err != nil {
	// 	panic(err)
	// }

	// logger := zap.Must(zap.NewProduction()).Sugar()

	// txManager := database.NewSQLTransactionManager(db)
	// testStore := store.NewPostgresStorage(db)

	// authServiceConfig := appconfig.GetAuthServiceConfig(configs)
	// testService := service.NewService(txManager, testStore, mailer, logger, jwtAuthenticator, tokenAuthenricator, authServiceConfig)

	// rateLimiter := appconfig.GetFixedWindowLimiterFromConfig(configs)

	// testApp := api.NewApp(
	// 	configs,
	// 	dataValidator,
	// 	testService,
	// 	logger,
	// 	rateLimiter,
	// )

	// testRouter := testApp.InitRouter() // Useful for http tests

	clearTestDB(db)
	seedUsersFuzz(f, db)

	f.Add("mario@example.com", "Password%123!")

	f.Fuzz(func(t *testing.T, email, password string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("PANIC: %v", r)
			}
		}()

		// Use some valid (generated with seeding) data
		if customRand.Float32() < 0.3 { // ~ 30% use valid email
			i := customRand.Intn(min(seedUserNum, len(validEmails), len(validPasswords)))

			email = validEmails[i]
			t.Log("Valid Email tested")

			if customRand.Float32() < 0.5 { // ~ 50% use valid password for the chosen email
				password = validPasswords[i]
				t.Log("Valid Password tested")
			}
		}

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
					t.Fatal("user not sent") // TODO: fix this error (response body: {"data":{"is_email_sent":false}})
				}

				if res.Data.User.Email == "" {
					t.Fatalf("empty email on success:\nReq: %+v\nRes:%+v", req, res)
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

// TODO: testa anche google e google/callback, verifiche varie (tutti gli handler dell'auth)
