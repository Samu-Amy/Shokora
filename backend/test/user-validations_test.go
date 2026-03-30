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

	t.Run("should allow empty optional fields", func(t *testing.T) {
		req := payloads.UserDataReq{
			FirstName: "Mario",
			LastName:  "",
			Birthday:  "",
		}

		err := dataValidator.Struct(req)

		if err != nil {
			t.Errorf("expected valid, got %v", err)
		}
	})

	t.Run("should not pass first_name validation", func(t *testing.T) {
		invalidField := "first_name"
		logErr := false

		for val, expectedTag := range notValidFirstNamesValidation {
			req := payloads.UserDataReq{
				FirstName: val,
				LastName:  randomFrom(validLastNames),
				Birthday:  randomFrom(validBirthdays),
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

			// Check validation
			found := false

			for _, ve := range validationErrors {
				if ve.Field() == invalidField && ve.Tag() == expectedTag {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error on field %s with tag %s not found", invalidField, expectedTag)
			}
		}
	})

	t.Run("should not pass last_name validation", func(t *testing.T) {
		invalidField := "last_name"
		logErr := false

		for val, expectedTag := range notValidLastNamesValidation {
			req := payloads.UserDataReq{
				FirstName: randomFrom(validFirstNames),
				LastName:  val,
				Birthday:  randomFrom(validBirthdays),
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

			// Check validation
			found := false

			for _, ve := range validationErrors {
				if ve.Field() == invalidField && ve.Tag() == expectedTag {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error on field %s with tag %s not found", invalidField, expectedTag)
			}
		}
	})

	t.Run("should not pass birthday validation", func(t *testing.T) {
		invalidField := "birthday"
		expectedTag := "valid-birthday"
		logErr := false

		for _, val := range notValidBirthdays {
			req := payloads.UserDataReq{
				FirstName: randomFrom(validFirstNames),
				LastName:  randomFrom(validLastNames),
				Birthday:  val,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Error("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

			// Check validation
			found := false

			for _, ve := range validationErrors {
				if ve.Field() == invalidField || ve.Tag() != expectedTag {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error on field %s with tag %s not found", invalidField, expectedTag)
			}
		}
	})

	t.Run("should not pass validation for more than one field", func(t *testing.T) {
		logErr := false

		req := payloads.UserDataReq{
			FirstName: "-Jènt&",
			LastName:  "ben#",
			Birthday:  "124-08",
		}

		err := dataValidator.Struct(req)

		if err == nil {
			t.Error("expected not valid, got valid")
		}

		if logErr {
			t.Logf("error: %v", err)
		}

		validationErrors := parseValidationErr(t, err)

		// Check validation
		foundFirstNameErr := false
		foundLastNameErr := false
		foundBirthdayErr := false

		for _, ve := range validationErrors {
			switch {
			case ve.Field() == "first_name" && ve.Tag() == "valid-name":
				foundFirstNameErr = true
			case ve.Field() == "last_name" && ve.Tag() == "valid-name":
				foundLastNameErr = true
			case ve.Field() == "birthday" && ve.Tag() == "valid-birthday":
				foundBirthdayErr = true
			}
		}

		if !(foundFirstNameErr && foundLastNameErr && foundBirthdayErr) {
			t.Error("expected error on all fields")
		}
	})
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

	t.Run("should not pass validation", func(t *testing.T) {
		invalidField := "email"
		logErr := false

		for val, expectedTag := range notValidEmailsValidation {
			req := payloads.EmailFieldReq{
				Email: val,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

			// Check validation
			found := false

			for _, ve := range validationErrors {
				if ve.Field() == invalidField && ve.Tag() == expectedTag {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error on field %s with tag %s not found", invalidField, expectedTag)
			}
		}
	})
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
