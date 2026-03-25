package payloads

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// - Regex -
var (
	regName     = regexp.MustCompile(`^[A-Za-zÀ-ÖØ-öø-ÿ' \-]+$`)
	regPassword = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()\-_=+\[\]{};:'",.<>?/\\|` + "`" + ` ~]+$`)
	regBirthday = regexp.MustCompile(`^(0[1-9]|[12]\d|3[01])-(0[1-9]|1[0-2])$`)
)

// - Validator -

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Set the field name to the one used in the json (for request (data) - response (error) json matching)
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0] // Get the first value of the json tag (the field name)

		if name == "-" { // Ignore field
			return ""
		}

		return name
	})

	// Register custom validation functions
	v.RegisterValidation("valid-name", validName)

	v.RegisterValidation("no-edge-spaces", noEdgeSpaces)
	v.RegisterValidation("valid-password", validPassword)
	v.RegisterValidation("no-common-password", noCommonPassword)

	v.RegisterValidation("valid-birthday", validBirthday)

	return v
}

// - Auth -

// Password
var commonPasswords = map[string]struct{}{
	"password":         {},
	"12345678":         {},
	"qwerty":           {},
	"admin":            {},
	"passwordpassword": {},
	"123456789012":     {},
	"qwertyuiop123":    {},
	"iloveyou123456":   {},
	"letmein1234567":   {},
	"welcome123456":    {},
	"admin12345678":    {},
	"monkey123456789":  {},
	"dragon12345678":   {},
	"shadow12345678":   {},
	"master12345678":   {},
	"football123456":   {},
	"princess123456":   {},
	"sunshine123456":   {},
	"baseball1234567":  {},
	"098765432109":     {},
	"123412341234":     {},
	"202312345678":     {},
	"202412345678":     {},
	"qwerty12345678":   {},
	"111111111111":     {},
	"999999999999":     {},
	"000000000000":     {},
	"654321654321":     {},
	"passwordpass123":  {},
	"passw0rdpass123":  {},
	"P@ssword123456":   {},
}

func IsCommonPassword(password string) bool {
	_, found := commonPasswords[strings.ToLower(password)]
	return found
}

func noCommonPassword(fl validator.FieldLevel) bool {
	return !IsCommonPassword(fl.Field().String())
}

func validName(fl validator.FieldLevel) bool {
	return regName.MatchString(fl.Field().String())
}

func noEdgeSpaces(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Do not accept spaces at start or end
	return strings.TrimSpace(value) == value
}

func validPassword(fl validator.FieldLevel) bool {
	return regPassword.MatchString(fl.Field().String())
}

func validBirthday(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Verify string format (regex)
	if !regBirthday.MatchString(value) {
		return false
	}

	// Verify date
	_, err := time.Parse("02-01-2006", value+"-2000") // 02-01-2006 means format the date as DD-MM-YYYY (and then using the value (with format "DD-MM") with "-2000" -> leap year, so 29-02 is valid)
	return err == nil
}
