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

func TestHashing3(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
