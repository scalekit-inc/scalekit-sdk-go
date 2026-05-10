---
operationId: RolesService_RemovePermissionFromRole
---

```javascript
await scalekit.permission.removePermissionFromRole("admin", "read:users");
```

```python
scalekit_client.permissions.remove_permission_from_role(
    role_name="admin",
    permission_name="read:users"
)
```

## Go SDK

```go
err := scalekitClient.Permission().RemovePermissionFromRole(ctx, "admin", "read:users")
if err != nil { /* handle err */ }
```

## Java SDK

```java
scalekitClient.permissions().removePermissionFromRole("admin", "read:users");
```
