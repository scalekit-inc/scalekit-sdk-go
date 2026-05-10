---
operationId: UserService_DeleteUser
---

```javascript
await scalekit.user.deleteUser("usr_123");
```

```python
scalekit_client.users.delete_user(organization_id="org_123", 
  user_id="usr_123")
```

## Go SDK

```go
if err := scalekitClient.User().DeleteUser(ctx, 
  "usr_123"); err != nil {
    panic(err)
}
```

## Java SDK

```java
users.deleteUser("usr_123");
```
