---
operationId: WebAuthnService_DeleteCredential
---

```javascript
const res = await scalekit.webauthn.deleteCredential("wac_123");
```

```python
res = scalekit_client.webauthn.delete_credential(credential_id="wac_123")
```

## Go SDK

```go
resp, err := scalekitClient.WebAuthn().DeleteCredential(ctx, "wac_123")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
DeleteCredentialResponse res = scalekitClient.webauthn().deleteCredential("wac_123");
```