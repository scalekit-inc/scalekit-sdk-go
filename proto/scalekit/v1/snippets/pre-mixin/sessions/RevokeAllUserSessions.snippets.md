---
operationId: SessionService_RevokeAllUserSessions
---

```javascript
const res = await scalekit.session.revokeAllUserSessions("user_123");
```

```python
res = scalekit_client.sessions.revoke_all_user_sessions(user_id="user_123")
```

## Go SDK

```go
resp, err := scalekitClient.Session().RevokeAllUserSessions(ctx, "user_123")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
RevokeAllUserSessionsResponse res = scalekitClient.sessions().revokeAllUserSessions("user_123");
```
