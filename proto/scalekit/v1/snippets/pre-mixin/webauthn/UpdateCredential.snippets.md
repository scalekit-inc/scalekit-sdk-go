---
operationId: WebAuthnService_UpdateCredential
---

```javascript
const res = await scalekit.webauthn.updateCredential(
  "wac_123",
  "Work Laptop Passkey"
);
```

```python
res = scalekit_client.webauthn.update_credential(
    credential_id="wac_123",
    display_name="Work Laptop Passkey"
)
```

## Go SDK

```go
resp, err := scalekitClient.WebAuthn().UpdateCredential(
    ctx,
    "wac_123",
    "Work Laptop Passkey",
)
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdateCredentialResponse res = scalekitClient.webauthn()
    .updateCredential("wac_123", "Work Laptop Passkey");
```