---
operationId: RolesService_CreatePermission
---

```javascript
await scalekit.permission.createPermission({
  name: "read:users",
  description: "Allows reading users"
});
```

```python
from scalekit.v1.roles.roles_pb2 import CreatePermission

permission = CreatePermission(
    name="read:users",
    description="Allows reading users"
)

scalekit_client.permissions.create_permission(permission=permission)

```

## Go SDK

```go
resp, err := scalekitClient.Permission().CreatePermission(ctx, &rolesv1.CreatePermission{
	Name:        "read:users",
	Description: "Allows reading users",
})
if err != nil { /* handle err */ }
_ = resp

```

## Java SDK

```java
CreatePermissionResponse res = scalekitClient.permissions().createPermission(
    CreatePermissionRequest.newBuilder()
        .setPermission(
            CreatePermission.newBuilder()
                .setName("read:users")
                .setDescription("Allows reading users")
                .build()
        )
        .build()
);

```
