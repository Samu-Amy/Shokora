package main

import (
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: controlla dati errore (field e tag, magari fare una lista di struct con field(s) values and name e tag dell'errore)

// User Data Req
func TestUserDataReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

		for range 50 {
			req := payloads.UserDataReq{
				FirstName: randomFrom(validFirstNames),
				LastName:  randomFrom(validLastNames),
				Birthday:  randomFrom(validBirthdays),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}

	})

	t.Run("should not pass first_name validation", func(t *testing.T) {

	})

	t.Run("should not pass last_name validation", func(t *testing.T) {

	})

	t.Run("should not pass birthday validation", func(t *testing.T) {

	})
}

// Email Field Req
func FuzzEmailFieldReq(f *testing.F) {
	f.Add("mario@example.com")

	f.Fuzz(func(t *testing.T, email string) {
		req := payloads.EmailFieldReq{Email: email}
		_ = dataValidator.Struct(req)
	})
}

func TestEmailFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

	})

	t.Run("should not pass validation", func(t *testing.T) {
	})
}

// - Password Field Req -

func FuzzPasswordFieldReq(f *testing.F) {
	f.Add("Password123!")

	f.Fuzz(func(t *testing.T, password string) {
		req := payloads.PasswordFieldReq{Password: password}
		_ = dataValidator.Struct(req)
	})
}

func TestPasswordFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

	})

	t.Run("should not pass validation", func(t *testing.T) {

	})
}

// - Double Password Field Req -

// func FuzzPasswordFieldReq(f *testing.F) {
// 	f.Add("Password123!")

// 	f.Fuzz(func(t *testing.T, oldPassword, newPassword string) {
// 		req := payloads.UpdatePasswordReq{
// 			OldPassword: oldPassword,
// 			NewPassword: newPassword,
// 		}
// 		_ = v.Struct(req)
// 	})
// }

func TestDoublePasswordFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

	})

	t.Run("should not pass validation", func(t *testing.T) {

	})
}

// - Update Password Field Req -

func TestUpdatePasswordReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

	})

	t.Run("should not pass validation", func(t *testing.T) {

	})
}
