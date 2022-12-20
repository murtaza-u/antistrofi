package token_test

import (
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/murtaza-u/antistrofi/token"
)

func TestNewToken(t *testing.T) {
	body := make(map[string]any)
	body["foo"] = "bar"

	var err error

	// invalid `nbf` field
	_, err = token.NewToken(token.Params{
		Body:      body,
		Expiry:    time.Now().Add(time.Minute * 5),
		NotBefore: time.Now().Add(time.Minute * 10),
	})
	if err != token.ErrInvalidNB {
		t.Errorf(
			"NewToken: expected: %s | got: %s",
			token.ErrInvalidNB.Error(), err.Error(),
		)
	}

	// body = nil
	_, err = token.NewToken(token.Params{
		Body: nil,
	})
	if err != token.ErrMissingBody {
		t.Errorf(
			"NewToken: expected: %s | got: %s",
			token.ErrMissingBody.Error(), err.Error(),
		)
	}

	// body = empty map
	_, err = token.NewToken(token.Params{
		Body: make(map[string]any),
	})
	if err != token.ErrMissingBody {
		t.Errorf(
			"NewToken: expected: %s | got: %s",
			token.ErrMissingBody.Error(), err.Error(),
		)
	}

	// no errors
	_, err = token.NewToken(token.Params{
		Body: body,
	})
	if err != nil {
		t.Errorf(
			"NewToken: expected: nil | got: %s", err.Error(),
		)
	}
}

func TestEncrypt(t *testing.T) {
	p := token.Params{
		Body: map[string]any{"foo": "bar"},
	}

	tkn, err := token.NewToken(p)
	if err != nil {
		t.Errorf("NewToken: %s", err.Error())
	}

	// invalid secret
	_, err = tkn.Encrypt([]byte("thiswillfail"), nil)
	if err == nil {
		t.Error("*token.Encrypt: expected: an error | got: nil")
	}

	// valid secret. Should pass.
	key := paseto.NewV4SymmetricKey()
	enc, err := tkn.Encrypt(key.ExportBytes(), nil)
	if err != nil {
		t.Errorf("*token.Encrypt: %s", err.Error())
	}
	t.Logf("encrypted token: %s\n", enc)
}

func TestIsExpired(t *testing.T) {
	p := token.Params{
		Expiry: time.Now().Add(time.Millisecond),
		Body:   map[string]any{"foo": "bar"},
	}

	tkn, err := token.NewToken(p)
	if err != nil {
		t.Errorf("NewToken: %s", err.Error())
	}

	time.Sleep(time.Millisecond)

	if !tkn.IsExpired() {
		t.Errorf("*token.IsExpired: invalid outcome")
	}
}

func TestRefresh(t *testing.T) {
	p := token.Params{
		Expiry: time.Now().Add(time.Millisecond),
		Body:   map[string]any{"foo": "bar"},
	}

	tkn, err := token.NewToken(p)
	if err != nil {
		t.Errorf("NewToken: %s", err.Error())
	}

	time.Sleep(time.Millisecond)

	err = tkn.Refresh()
	if err != nil {
		t.Errorf("Refresh: %s", err.Error())
	}
}
