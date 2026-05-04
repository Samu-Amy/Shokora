package main

import "testing"

func assertValidationFails(t *testing.T, err error, expectedField, expectedTag, val string) {
	t.Helper()

	if err == nil {
		if val == "" {
			t.Fatal("expected validation error, got nil")
		} else {
			t.Fatalf("expected validation error, got nil on val: %q", val)
		}
	}

	vErrs := parseValidationErr(t, err)

	for _, ve := range vErrs {
		if ve.Field() == expectedField && ve.Tag() == expectedTag {
			return
		}
	}

	if val != "" {
		t.Errorf("expected error on field %q with tag %q, got: %v on val: %q", expectedField, expectedTag, vErrs, val)
	} else {
		t.Errorf("expected error on field %q with tag %q, got: %v", expectedField, expectedTag, vErrs)
	}
}

func assertValidationFailsWithParam(t *testing.T, err error, expectedField, expectedTag, expectedParam, val string) {
	t.Helper()

	if err == nil {
		if val == "" {
			t.Fatal("expected validation error, got nil")
		} else {
			t.Fatalf("expected validation error, got nil on val: %q", val)
		}
	}

	vErrs := parseValidationErr(t, err)

	for _, ve := range vErrs {
		if ve.Field() == expectedField && ve.Tag() == expectedTag && ve.Param() == expectedParam {
			return
		}
	}

	if val != "" {
		t.Errorf("expected error on field %q with tag %q and param %q, got: %v on val: %q", expectedField, expectedTag, expectedParam, vErrs, val)
	} else {
		t.Errorf("expected error on field %q with tag %q and param %q, got: %v", expectedField, expectedTag, expectedParam, vErrs)
	}
}
