package auth

import "errors"

// ValidatePasswordStrength does a quick sanity check before the web app
// even calls the API's register endpoint, so the user gets instant
// feedback instead of a round-trip for something like an empty password.
// The API is still the source of truth for hashing (bcrypt) and storage.
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
