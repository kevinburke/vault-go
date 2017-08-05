package vault

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTransit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	c := NewClient("http://localhost:8200", "8928a0a0-0584-7d10-8fe7-de6c76d2e685", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := c.Transit.Create(ctx, "test-key", map[string]interface{}{
		"type": "ed25519",
	})
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("Alas, poor Horatio")
	sig, err := c.Transit.Sign(ctx, "test-key", msg, map[string]interface{}{
		"algorithm": "ed25519",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sig.Signature)
	valid, err := c.Transit.VerifySignature(ctx, "test-key", msg, sig.Signature, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Errorf("expected to validate signature, got false")
	}
}
