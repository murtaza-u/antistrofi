package token

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

// GetFooter extracts the footer from the encrypted token. No secret key
// is required since footer is never encrypted.
func GetFooter(enc string) ([]byte, error) {
	p := paseto.NewParser()
	footer, err := p.UnsafeParseFooter(protocol, enc)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to extract footer: %s", err.Error(),
		)
	}
	return footer, nil
}

// Decrypt decrypts the encrypted token
func Decrypt(enc string, secret, implicit []byte) (*token, error) {
	key, err := paseto.V4SymmetricKeyFromBytes(secret)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create symmetric key from provided secret: %s",
			err.Error(),
		)
	}

	p := paseto.NewParserWithoutExpiryCheck()
	tkn, err := p.ParseV4Local(key, enc, implicit)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt token: %s", err.Error(),
		)
	}

	return &token{Token: *tkn}, nil
}
