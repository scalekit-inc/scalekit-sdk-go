---
operationId: RolesService_UpdateRole
---

```javascript
await scalekit.role.updateRole("admin", {
  displayName: "Admin (Updated)",
  description: "Updated description"
});
```

```python
from scalekit.v1.roles.roles_pb2 import UpdateRole

scalekit_client.roles.update_role(
    role_name="admin",
    role=UpdateRole(
        display_name="Admin (Updated)",
        description="Updated description"
    )
)
```

## Go SDK

```go
resp, err := scalekitClient.Role().UpdateRole(ctx, "admin", &rolesv1.UpdateRole{
    DisplayName: "Admin (Updated)",
    Description: "Updated description",
})
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdateRoleResponse res = scalekitClient.roles().updateRole(
    "admin",
    UpdateRoleRequest.newBuilder()
        .setRole(
            UpdateRole.newBuilder()
                .setDisplayName("Admin (Updated)")
                .setDescription("Updated description")
                .build()
        )
        .build()
);
```
