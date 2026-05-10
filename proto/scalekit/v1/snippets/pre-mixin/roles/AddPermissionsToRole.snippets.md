---
operationId: RolesService_AddPermissionsToRole
---

```javascript
await scalekit.permission.addPermissionsToRole("role_admin", ["perm.read", "perm.write"]);
```

```python
scalekit_client.permissions.add_permissions_to_role(
    role_name="role_admin",
    permission_names=["perm.read", "perm.write"]
)
```

## Go SDK

```go
resp, err := scalekitClient.Permission().AddPermissionsToRole(ctx, "role_admin", []string{"perm.read", "perm.write"})

```

## Java SDK

```java
AddPermissionsToRoleResponse res = scalekitClient.permissions().addPermissionsToRole(
    "role_admin",
    AddPermissionsToRoleRequest.newBuilder()
        .addPermissionNames("perm.read")
        .addPermissionNames("perm.write")
        .build()
);
```
