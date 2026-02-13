package validator

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Valid", "valid_user_123", false},
		{"Short", "ab", true},
		{"Long", "this_username_is_way_too_long_and_should_fail_validation_because_it_exceeds_50_chars", true},
		{"InvalidChar", "user@name", true},
		{"Space", "user name", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateUsername(tt.username); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid", "StrongP@ssw0rd!", false},
		{"Short", "Weak1!", true},
		{"NoUpper", "weakpassword1!", true},
		{"NoDigit", "WeakPassword!", true},
		{"NoSpecial", "WeakPassword1", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePassword(tt.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Plain", "Hello World", "Hello World"},
		{"Script", "<script>alert('xss')</script>", ""},
		{"Tags", "<b>Bold</b>", "Bold"},
		{"AllowedTags", "<b>Hello</b>", "Hello"},
		{"Dangerous", "<a href='javascript:alert(1)'>Link</a>", "Link"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sanitize(tt.input)
			if got != tt.expected {
				t.Errorf("Sanitize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
