package vault

import (
	"context"
	"encoding/base64"
)

const transitKeyPathPart = "/transit/keys"
const transitSignPathPart = "/transit/sign"
const transitVerifyPathPart = "/transit/verify"

// TransitService contains helpers for working with the Transit API.
type TransitService struct {
	client *Client
}

// Create creates a new Key with the specified data. For more on valid inputs
// see https://www.vaultproject.io/api/secret/transit/index.html#create-key.
func (t *TransitService) Create(ctx context.Context, name string, data map[string]interface{}) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	return t.client.UpdateResource(ctx, transitKeyPathPart, name, data, nil)
}

type Signature struct {
	Signature string `json:"signature"`
}

// Sign creates a new signature for the input using the secret
// key. The input should not be base64 encoded. For more, see
// https://www.vaultproject.io/api/secret/transit/index.html#sign-data.
func (t *TransitService) Sign(ctx context.Context, key string, input []byte, data map[string]interface{}) (*Signature, error) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["input"] = base64.StdEncoding.EncodeToString(input)
	sig := new(Signature)
	r := &Response{
		Data: sig,
	}
	err := t.client.UpdateResource(ctx, transitSignPathPart, key, data, r)
	return sig, err
}

type valid struct {
	Valid bool `json:"valid"`
}

// Verify verifies the input using the secret key. input should not be base64
// encoded. For more see
// https://www.vaultproject.io/api/secret/transit/index.html#verify-signed-data
func (t *TransitService) VerifySignature(ctx context.Context, key string, input []byte, sig string, data map[string]interface{}) (bool, error) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["signature"] = sig
	data["input"] = base64.StdEncoding.EncodeToString(input)
	v := new(valid)
	r := &Response{
		Data: v,
	}
	err := t.client.UpdateResource(ctx, transitVerifyPathPart, key, data, r)
	return v.Valid, err
}
