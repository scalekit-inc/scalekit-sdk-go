---
operationId: UserService_ListUserRoles
---

```javascript
const { roles } = await scalekit.user.listUserRoles("org_123", "usr_123");
```

```python
resp = scalekit_client.users.list_user_roles(
    organization_id="org_123",
    user_id="usr_123",
)
roles = resp.roles
```

```go
resp, err := scalekitClient.User().ListUserRoles(ctx, "org_123", "usr_123")
if err != nil {
    // handle error
}
roles := resp.Roles
```

```java
ListUserRolesResponse resp = scalekitClient.users().listUserRoles("org_123", "usr_123");
List<Role> roles = resp.getRolesList();
```
