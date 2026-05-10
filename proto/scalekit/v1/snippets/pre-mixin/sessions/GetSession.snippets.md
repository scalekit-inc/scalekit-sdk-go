---
operationId: SessionService_GetSession
---

```javascript
const res = await scalekit.session.getSession("ses_123456789");
```

```python
res = scalekit_client.sessions.get_session(session_id="ses_123456789")
```

## Go SDK

```go
resp, err := scalekitClient.Session().GetSession(ctx, "ses_123456789")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
SessionDetails res = scalekitClient.sessions().getSession("ses_123456789");
```
