package vault_test

import (
	"context"
	"fmt"
	"log"

	vault "github.com/kevinburke/vault-go"
)

func ExampleTransitService_Create() {
	c := vault.NewClient("https://localhost:8200", "your-token", nil)
	data := map[string]interface{}{
		"exportable": false,
		"algorithm":  "ed25519",
		"derived":    false,
	}
	err := c.Transit.Create(context.TODO(), "my-key", data)
	fmt.Println(err)
}

func ExampleTransitService_roundtrip() {
	c := vault.NewClient("https://localhost:8200", "your-token", nil)
	msg := []byte("Alas, poor Yorick")
	sig, err := c.Transit.Sign(context.TODO(), "my-key", msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	valid, err := c.Transit.VerifySignature(context.TODO(), "my-key", msg, sig.Signature, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("valid", valid)
}
