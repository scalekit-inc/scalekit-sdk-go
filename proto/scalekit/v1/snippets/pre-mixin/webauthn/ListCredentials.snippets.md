---
operationId: WebAuthnService_ListCredentials
---

```javascript
const res = await scalekit.webauthn.listCredentials("user_123");
```

```python
res = scalekit_client.webauthn.list_credentials(user_id="user_123")
```

```go
resp, err := scalekitClient.WebAuthn().ListCredentials(ctx, "user_123")
if err != nil { /* handle err */ }
_ = resp
```

```java
ListCredentialsResponse res = scalekitClient.webauthn().listCredentials("user_123");
```