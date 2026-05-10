---
operationId: UsersService_ResendInvite
---

```javascript
const response = await scalekit.user.resendInvite("org_123", "usr_456");
// Invitation email resent to user
```

```python
response = scalekit_client.user.resend_invite(
    organization_id="org_123",
    user_id="usr_456"
)
# Invitation email resent to user
```

```go
resp, err := scalekitClient.User().ResendInvite(ctx, "org_123", "usr_456")
if err != nil {
    // handle error
}
// Invitation email resent to user
```

```java
ResendInviteResponse resp = scalekitClient.users().resendInvite("org_123", "usr_456");
// Invitation email resent to user
```
