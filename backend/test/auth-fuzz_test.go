package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
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

func FuzzLoginUserRoute(f *testing.F) { // go test .\test\ -run=^$ -fuzz=FuzzLoginUserRoute -fuzztime=20s

	clearTestDB(db)
	seedUsersFuzz(f, db)

	f.Add("mario@example.com", "Password%123!")

	f.Fuzz(func(t *testing.T, email, password string) {

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

		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", req)

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
			if res.Data.User.Email == "" {
				t.Fatalf("empty email on success:\nReq: %+v\nRes:%+v", req, res)
			}

			if res.Data.User.IsVerified || res.Data.User.Role != 0 {
				t.Fatal("user created as verified or with role != 0")
			}
		}
	})
}

// TODO: testa anche google e google/callback, verifiche varie (tutti gli handler dell'auth)
