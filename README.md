Package vault is a client around the Vault v1 API.

Operations center around a Client. Create a client by specifying a base URL,
a token, and (optionally) a http.Client instance:

 client := vault.NewClient("https://localhost:8200", "my-token", nil)

Once you have a client, you can perform various operations.

    client.Transit.Create(context.TODO(), "key", map[string]interface{}{
       "algorithm": "ed25519",
    })

POST operations take all required parameters as required and optional
parameters as map[string]interface{} inputs. Consult the API documentation to
determine what optional parameters are available for each endpoint.
