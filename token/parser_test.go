package token_test

import (
	"testing"

	"aidanwoods.dev/go-paseto"
	"github.com/murtaza-u/antistrofi/token"
)

func TestGetFooter(t *testing.T) {
	tkn, err := token.New(token.Params{
		Body:   map[string]any{"foo": "bar"},
		Footer: "blah",
	})
	if err != nil {
		t.Errorf("NewToken: %s", err.Error())
	}

	key := paseto.NewV4SymmetricKey()
	enc, err := tkn.Encrypt(key.ExportBytes(), nil)
	if err != nil {
		t.Errorf("*token.Encrypt: %s", err.Error())
	}

	footer, err := token.GetFooter(enc)
	if err != nil {
		t.Errorf("GetFooter: %s", err.Error())
	}

	if string(footer) != "blah" {
		t.Errorf("incorrectly extracted foooter")
	}
}

func TestDecrypt(t *testing.T) {
	tkn, err := token.New(token.Params{
		Body: map[string]any{"foo": "bar"},
	})
	if err != nil {
		t.Errorf("NewToken: %s", err.Error())
	}

	key := paseto.NewV4SymmetricKey()
	enc, err := tkn.Encrypt(key.ExportBytes(), nil)
	if err != nil {
		t.Errorf("*token.Encrypt: %s", err.Error())
	}

	tkn, err = token.Decrypt(enc, key.ExportBytes(), nil)
	if err != nil {
		t.Errorf("*token.Decrypt: %s", err.Error())
	}

	bar, err := tkn.GetString("foo")
	if bar != "bar" {
		t.Error("incorrectly decrypted token")
	}
}
