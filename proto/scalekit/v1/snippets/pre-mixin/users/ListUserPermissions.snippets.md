---
operationId: UserService_ListUserPermissions
---

```javascript
const { permissions } = await scalekit.user.listUserPermissions("org_123", "usr_123");
```

```python
resp = scalekit_client.users.list_user_permissions(
    organization_id="org_123",
    user_id="usr_123",
)
permissions = resp.permissions
```

```go
resp, err := scalekitClient.User().ListUserPermissions(ctx, "org_123", "usr_123")
if err != nil {
    // handle error
}
permissions := resp.Permissions
```

```java
ListUserPermissionsResponse resp = scalekitClient.users().listUserPermissions("org_123", "usr_123");
List<Permission> permissions = resp.getPermissionsList();
```
