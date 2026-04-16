package main

import (
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/google/uuid"
)

// TODO: fix (controlla sia test che mock data)

// Google OAuth Callback Req
func TestGoogleOAuthCallbackReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.GoogleOAuthCallbackReq{
				State: randomFrom(validBase64RawUrl32BytesString),
				Code:  randomFrom(validGoogleCodes),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass state validation", func(t *testing.T) {
		invalidField := "state"
		logErr := false

		for val, expectedTag := range notValidBase64RawUrl32BytesString {
			req := payloads.GoogleOAuthCallbackReq{
				State: val,
				Code:  randomFrom(validGoogleCodes),
			}

			err := dataValidator.Struct(req)

			assertValidationFails(t, err, invalidField, expectedTag, val)

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}
		}
	})

	t.Run("should not pass code validation", func(t *testing.T) {
		invalidField := "code"
		logErr := false

		for val, expectedTag := range notValidGoogleCodesValidation {
			req := payloads.GoogleOAuthCallbackReq{
				State: randomFrom(validBase64RawUrl32BytesString),
				Code:  val,
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

	t.Run("should not pass validation for both fields", func(t *testing.T) {
		logErr := false

		req := payloads.GoogleOAuthCallbackReq{
			State: "invalid state!@#$",
			Code:  "invalid@code!#$%",
		}

		err := dataValidator.Struct(req)

		if err == nil {
			t.Error("expected not valid, got valid")
		}

		if logErr {
			t.Logf("error: %v", err)
		}

		validationErrors := parseValidationErr(t, err)

		foundStateErr := false
		foundCodeErr := false

		for _, ve := range validationErrors {
			switch {
			case ve.Field() == "state" && ve.Tag() == "safe-chars":
				foundStateErr = true
			case ve.Field() == "code" && ve.Tag() == "safe-chars":
				foundCodeErr = true
			}
		}

		if !(foundStateErr && foundCodeErr) {
			t.Error("expected error on both state and code fields")
		}
	})
}

// OTP Verification Req
func TestOTPVerificationReqValidation(t *testing.T) {
	validUUID := uuid.New()

	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.OTPVerificationReq{
				VerificationId: validUUID,
				OTP:            randomFrom(validOTPs),
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass verification_id validation for nil UUID", func(t *testing.T) {
		invalidField := "verification_id"
		expectedTag := "required"
		logErr := false

		req := payloads.OTPVerificationReq{
			VerificationId: uuid.Nil,
			OTP:            randomFrom(validOTPs),
		}

		err := dataValidator.Struct(req)

		if err == nil {
			t.Fatal("expected not valid, got valid")
		}

		if logErr {
			t.Logf("error: %v", err)
		}

		validationErrors := parseValidationErr(t, err)

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
	})

	t.Run("should not pass OTP validation", func(t *testing.T) {
		invalidField := "otp"
		logErr := false

		for val, expectedTag := range notValidOTPsValidation {
			req := payloads.OTPVerificationReq{
				VerificationId: validUUID,
				OTP:            val,
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

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

	t.Run("should not pass validation for both fields", func(t *testing.T) {
		invalidField1 := "verification_id"
		invalidField2 := "otp"
		logErr := false

		req := payloads.OTPVerificationReq{
			VerificationId: uuid.Nil,
			OTP:            "invalid",
		}

		err := dataValidator.Struct(req)

		if err == nil {
			t.Error("expected not valid, got valid")
		}

		if logErr {
			t.Logf("error: %v", err)
		}

		validationErrors := parseValidationErr(t, err)

		foundVerificationIDErr := false
		foundOTPErr := false

		for _, ve := range validationErrors {
			if ve.Field() == invalidField1 && ve.Tag() == "required" {
				foundVerificationIDErr = true
			} else if ve.Field() == invalidField2 && ve.Tag() == "valid-otp" {
				foundOTPErr = true
			}
		}

		if !(foundVerificationIDErr && foundOTPErr) {
			t.Error("expected error on both verification_id and otp fields")
		}
	})
}

// Reset Password Req
func TestResetPasswordReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.ResetPasswordReq{
				PlainResetSessionToken: randomFrom(validBase64RawUrl32BytesString),
				PasswordFieldReq: payloads.PasswordFieldReq{
					Password: randomFrom(validPasswords),
				},
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass plain_reset_session_token validation", func(t *testing.T) {
		invalidField := "plain_reset_session_token"
		logErr := false

		for val, expectedTag := range notValidBase64RawUrl32BytesString {
			req := payloads.ResetPasswordReq{
				PlainResetSessionToken: val,
				PasswordFieldReq: payloads.PasswordFieldReq{
					Password: randomFrom(validPasswords),
				},
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

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

	t.Run("should not pass password validation", func(t *testing.T) {
		invalidField := "password"
		logErr := false

		for val, expectedTag := range notValidPasswordsValidation {
			req := payloads.ResetPasswordReq{
				PlainResetSessionToken: randomFrom(validBase64RawUrl32BytesString),
				PasswordFieldReq: payloads.PasswordFieldReq{
					Password: val,
				},
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

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

	t.Run("should not pass validation for both fields", func(t *testing.T) {
		logErr := false

		req := payloads.ResetPasswordReq{
			PlainResetSessionToken: "invalid state!@#$",
			PasswordFieldReq: payloads.PasswordFieldReq{
				Password: "short",
			},
		}

		err := dataValidator.Struct(req)

		if err == nil {
			t.Error("expected not valid, got valid")
		}

		if logErr {
			t.Logf("error: %v", err)
		}

		validationErrors := parseValidationErr(t, err)

		foundTokenErr := false
		foundPasswordErr := false

		for _, ve := range validationErrors {
			switch {
			case ve.Field() == "plain_reset_session_token" && ve.Tag() == "safe-chars":
				foundTokenErr = true
			case ve.Field() == "password" && ve.Tag() == "min":
				foundPasswordErr = true
			}
		}

		if !(foundTokenErr && foundPasswordErr) {
			t.Error("expected error on both plain_reset_session_token and password fields")
		}
	})
}

// Send Verification Req
func TestSendVerificationReqValidation(t *testing.T) {
	t.Run("should pass validation", func(t *testing.T) {
		for range validationTestsNum {
			req := payloads.SendVerificationReq{
				EmailFieldReq: payloads.EmailFieldReq{
					Email: randomFrom(validEmails),
				},
			}

			err := dataValidator.Struct(req)

			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		}
	})

	t.Run("should not pass email validation", func(t *testing.T) {
		invalidField := "email"
		logErr := false

		for val, expectedTag := range notValidEmailsValidation {
			req := payloads.SendVerificationReq{
				EmailFieldReq: payloads.EmailFieldReq{
					Email: val,
				},
			}

			err := dataValidator.Struct(req)

			if err == nil {
				t.Fatal("expected not valid, got valid")
			}

			if logErr {
				t.Logf("val: %s, error: %v", val, err)
			}

			validationErrors := parseValidationErr(t, err)

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
