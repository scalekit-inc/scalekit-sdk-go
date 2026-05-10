---
operationId: SessionService_RevokeSession
---

```javascript
const res = await scalekit.session.revokeSession("ses_123456789");
```

```python
res = scalekit_client.sessions.revoke_session(session_id="ses_123456789")
```

## Go SDK

```go
resp, err := scalekitClient.Session().RevokeSession(ctx, "ses_123456789")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
RevokeSessionResponse res = scalekitClient.sessions().revokeSession("ses_123456789");
```
