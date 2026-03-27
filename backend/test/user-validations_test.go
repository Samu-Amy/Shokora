package main

import (
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: controlla dati errore (field e tag, magari fare una lista di struct con field(s) values and name e tag dell'errore)

// User Data Req
func TestUserDataReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {

		for range validationTestsNum {
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

	// t.Run("should not pass first_name validation", func(t *testing.T) {

	// })

	// t.Run("should not pass last_name validation", func(t *testing.T) {

	// })

	// t.Run("should not pass birthday validation", func(t *testing.T) {

	// })
}

// Email Field Req
func TestEmailFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.EmailFieldReq{
				Email: randomFrom(validEmails),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	// t.Run("should not pass validation", func(t *testing.T) {

	// })
}

// - Password Field Req -

func TestPasswordFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.PasswordFieldReq{
				Password: randomFrom(validPasswords),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	// t.Run("should not pass validation", func(t *testing.T) {

	// })
}

// - Double Password Field Req -

func TestDoublePasswordFieldReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			passw := randomFrom(validPasswords)

			req := payloads.DoublePasswordFieldReq{
				Password:             passw,
				PasswordConfirmation: passw,
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass validation because of different password", func(t *testing.T) {
		for range validationTestsNum {
			passw1 := randomFrom(validPasswords)
			passw2 := randomFrom(validPasswords)

			for passw1 == passw2 {
				passw2 = randomFrom(validPasswords)
			}

			req := payloads.DoublePasswordFieldReq{
				Password:             passw1,
				PasswordConfirmation: passw2,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Errorf("expected invalid")
			}
		}
	})

	t.Run("should not pass validation because of ...", func(t *testing.T) {
		for range validationTestsNum {
			passw1 := randomFrom(validPasswords)
			passw2 := randomFrom(validPasswords)

			for passw1 == passw2 {
				passw2 = randomFrom(validPasswords)
			}

			req := payloads.DoublePasswordFieldReq{
				Password:             passw1,
				PasswordConfirmation: passw2,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Errorf("expected invalid")
			}
		}
	})
}

// - Update Password Field Req -

func TestUpdatePasswordReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.UpdatePasswordReq{
				OldPassword: randomFrom(validPasswords),
				NewPassword: randomFrom(validPasswords),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass validation bacause of same password", func(t *testing.T) {
		for range validationTestsNum {
			passw := randomFrom(validPasswords)

			req := payloads.DoublePasswordFieldReq{
				Password:             passw,
				PasswordConfirmation: passw,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Errorf("expected invalid")
			}
		}
	})
}
