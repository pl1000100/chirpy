package auth

import "testing"

func TestHashing1(t *testing.T) {
	h, _ := HashPassword("password")
	result := CheckPasswordHash(h, "password")
	var expected error = nil

	if result != expected {
		t.Errorf("CheckPasswordHash(HashPassword(\"password\"), \"password\") = %d; want %d", result, expected)
	}
}

func TestHashing2(t *testing.T) {
	h, _ := HashPassword("password")
	result := CheckPasswordHash(h, "password1")
	var expected error = nil

	if result == expected {
		t.Errorf("CheckPasswordHash(HashPassword(\"password\"), \"password\") = %d; don't want %d", result, expected)
	}
}
