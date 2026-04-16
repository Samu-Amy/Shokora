package main

// - OTP -

var validOTPs = []string{
	"000000", "000001", "123456", "999999", "555555", "111111", "666666", "777777", "888888", "222222",
	"000000", "111111", "222222", "333333", "444444", "555555", "666666", "777777", "888888", "999999",
	"100000", "200000", "300000", "400000", "500000", "600000", "700000", "800000", "900000", "050000",
	"763296", "673296", "198329", "289173", "864327", "098412", "542134", "986759", "903272", "874236",
}

var notValidOTPsValidation = map[string]string{
	"":        "required",  // vuoto
	"12345":   "valid-otp", // meno di 6 cifre
	"1234567": "valid-otp", // più di 6 cifre
	"12345a":  "valid-otp", // contiene lettere
	"123 456": "valid-otp", // contiene spazio
	"123-456": "valid-otp", // contiene simbolo
	"abcdef":  "valid-otp", // solo lettere
	"123@56":  "valid-otp", // contiene simbolo
}

// - Google OAuth State and Magic Link / Reset Password Token -

var validBase64RawUrl32BytesString = []string{
	// Token da 32 byte convertiti a RawURLBase64 (43 chars)
	"tzb6cR1saaiYuhwlbqvwTz0RGx9CkBaueirSqCKl8S5", // 43
	"hZ086yr410hX6SmganhW0T32iyfKebXcCAm38cggxF0", // 43
	"Or71zHGPKEDE89eOyxiWZwB0yUlyC12Uoz9Xfzat3PM", // 43
	"Vz7qSSlTWsM3mf86IJ_kC5KEq93wDXnRSdBx-1JGnXw", // 43
	"82rfjnA3a_8Y-ykfulCAo2Brg-MFnMnb2lzpjs2ZbyQ", // 43
	"JfpnoS08ABPE87Lsa8_0eXCBbH3SmUIoW74wsFr3g7f", // 43
	"TeiLKXyTi-boyntxhSPdv_QWvDJbdWEesYrtUWWO6zw", // 43
	"7bGLXe4bCrdxLz4j8l08hELlMqJEo1CfnXSNsgaSZTg", // 43
	"3xo7m5CUaLun9LkrrvwY_OA-m8Oe5v1HRL-trq4NqaM", // 43
	"JbVEDTAvlyQ-o3IrPRCwh57FsbSNNcVKqdzgWigTPHk", // 43
}

var notValidBase64RawUrl32BytesString = map[string]string{
	"":      "required",
	"short": "valid-base64-rawurl-32", // Short
	"dGVzdHRva2VuMTIzNDU2Nzg5MDEyMzQ1Njc4OTA===":   "valid-base64-rawurl-32", // Short
	"dGVzdHRva2VuMTIzNDU2Nzg5MDEyMzQ1Njc4OTA87j":   "valid-base64-rawurl-32", // Short
	"dGVzdHRva2VuMTIzNDU2Nzg5MDEyMzQ1Njc4OTA87j=":  "valid-base64-rawurl-32", // Right length but padding ('=')
	"dGVzdHRva2VuMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM=": "valid-base64-rawurl-32", // Long
	"dGVzdHRva2VuMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjMs": "valid-base64-rawurl-32", // Long
	"invalidstate!@#$M😎IzNDU2Nzg5MDEyNjc4OTAxMjM":  "valid-base64-rawurl-32", // Char not val (Emojii)
	"invalidstate!@#$MIzNDU2Nzg5©MDEyNjc4OTAxMjM":  "valid-base64-rawurl-32", // Char not val ('©')
	"state with spaces HRva2VuMTIzNDU2Nzg5EyNjc4":  "valid-base64-rawurl-32", // Spaces
}

// - Google OAuth Code -

var validGoogleCodes = []string{
	"4/0AY0e-g4kQ",
	"4/0AX4XfQa",
	"4/0AY0e-g5Z",
	"somevalidgooglecode",
	"4_0AY0e_g4kQ",
	"gAAAAABmJ",
	"codewithlettersandnumbers123456",
	"4-0-AY0e_g4kQ",
	"validtokenwithunderscores_and_hyphens",
	"abc123def456ghi789jkl",
}

var notValidGoogleCodesValidation = map[string]string{
	"":                 "required",   // vuoto
	"invalid@code!#$%": "safe-chars", // caratteri non validi
	// "code with spaces":  "", // contiene spazi // TODO: non deve contenere spazi?
	"code\nwithnewline": "safe-chars", // newline
	"invalid©code":      "safe-chars", // carattere non valido
}

// ----- FUNCTIONS FOR VERIFICATION -----

// func makeGoogleOAuthCallbackReq(state, code string) payloads.GoogleOAuthCallbackReq {
// 	return payloads.GoogleOAuthCallbackReq{
// 		State: state,
// 		Code:  code,
// 	}
// }

// func makeOTPVerificationReq(verificationID string, otp string) payloads.OTPVerificationReq {
// 	uid, _ := uuid.Parse(verificationID) // Per i test, possiamo usare un UUID valido
// 	return payloads.OTPVerificationReq{
// 		VerificationId: uid,
// 		OTP:            otp,
// 	}
// }

// func makeSendVerificationReq(email string) payloads.SendVerificationReq {
// 	return payloads.SendVerificationReq{
// 		EmailFieldReq: payloads.EmailFieldReq{
// 			Email: email,
// 		},
// 	}
// }

// func makeResetPasswordReq(token, password string) payloads.ResetPasswordReq {
// 	return payloads.ResetPasswordReq{
// 		PlainResetSessionToken: token,
// 		PasswordFieldReq: payloads.PasswordFieldReq{
// 			Password: password,
// 		},
// 	}
// }
