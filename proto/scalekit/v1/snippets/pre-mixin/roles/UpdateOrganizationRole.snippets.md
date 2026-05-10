---
operationId: RolesService_UpdateOrganizationRole
---

```javascript
await scalekit.role.updateOrganizationRole("org_123", "org_admin", {
  displayName: "Org Admin (Updated)",
  description: "Updated org role description"
});
```

```python
from scalekit.v1.roles.roles_pb2 import UpdateRole

scalekit_client.roles.update_organization_role(
    org_id="org_123",
    role_name="org_admin",
    role=UpdateRole(
        display_name="Org Admin (Updated)",
        description="Updated org role description"
    )
)
```

## Go SDK

```go
resp, err := scalekitClient.Role().UpdateOrganizationRole(ctx, "org_123", "org_admin", &rolesv1.UpdateRole{
    DisplayName: "Org Admin (Updated)",
    Description: "Updated org role description",
})
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdateOrganizationRoleResponse res = scalekitClient.roles().updateOrganizationRole(
    "org_123",
    "org_admin",
    UpdateOrganizationRoleRequest.newBuilder()
        .setRole(
            UpdateRole.newBuilder()
                .setDisplayName("Org Admin (Updated)")
                .setDescription("Updated org role description")
                .build()
        )
        .build()
);
```
