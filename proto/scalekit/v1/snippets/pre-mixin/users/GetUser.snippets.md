---
operationId: UserService_GetUser
---

```javascript
const { user } = await scalekit.user.getUser("usr_123456");
```

```python
resp = scalekit_client.users.get_user(user_id="usr_123456")
user = resp.user
```

```go
resp, err := scalekitClient.User().GetUser(ctx, "usr_123456")
if err != nil {
    // handle error
}
user := resp.User
```

```java
GetUserResponse resp = scalekitClient.users().getUser("usr_123456");
User user = resp.getUser();
```
