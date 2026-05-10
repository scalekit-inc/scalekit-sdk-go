---
operationId: RolesService_UpdatePermission
---

```javascript
await scalekit.permission.updatePermission("read:users", {
  name: "read:users",
  description: "Allows reading user resources"
});
```

```python
from scalekit.v1.roles.roles_pb2 import CreatePermission

scalekit_client.permissions.update_permission(
    permission_name="read:users",
    permission=CreatePermission(
        name="read:users",
        description="Allows reading user resources"
    )
)
```

## Go SDK

```go
resp, err := scalekitClient.Permission().UpdatePermission(ctx, "read:users", &rolesv1.CreatePermission{
    Name:        "read:users",
    Description: "Allows reading user resources",
})
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdatePermissionResponse res = scalekitClient.permissions().updatePermission(
    "read:users",
    UpdatePermissionRequest.newBuilder()
        .setPermission(
            CreatePermission.newBuilder()
                .setName("read:users")
                .setDescription("Allows reading user resources")
                .build()
        )
        .build()
);
```
