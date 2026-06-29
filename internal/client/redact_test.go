package client

import (
	"strings"
	"testing"
)

func TestRedact(t *testing.T) {
	t.Parallel()

	input := `Authorization: Bearer abc123 Authorization: rawtoken {"raw_password":"secret-value","token":"token-value","password":"hash-value","reply_body":"private","smtp":"smtp://user:relaypass@mail.example.com"} smtp://user:relaypass@mail.example.com $bcrypt-sha256$v=2,t=2b,r=12$abc$def`
	output := Redact(input)

	for _, secret := range []string{"abc123", "rawtoken", "secret-value", "token-value", "hash-value", "private", "relaypass", "$bcrypt-sha256"} {
		if strings.Contains(output, secret) {
			t.Fatalf("redacted output still contains %q: %s", secret, output)
		}
	}
	if !strings.Contains(output, "<redacted>") {
		t.Fatalf("redacted output missing marker: %s", output)
	}
}
