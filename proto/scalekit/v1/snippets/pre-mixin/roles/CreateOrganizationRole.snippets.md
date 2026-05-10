---
operationId: RolesService_CreateOrganizationRole
---

```javascript
await scalekit.role.createOrganizationRole("org_123", {
  name: "org_admin",
  displayName: "Org Admin",
  description: "Organization-scoped role",
   extends: "base_role_name", // optional
   permissions: ["perm.read", "perm.write"] // optional
});
```

```python
from scalekit.v1.roles.roles_pb2 import CreateOrganizationRole

role = CreateOrganizationRole(
    name="org_admin",
    display_name="Org Admin",
    description="Organization-scoped role",
    extends="base_role_name",              # optional
    permissions=["perm.read", "perm.write"]  # optional
)

scalekit_client.roles.create_organization_role(
    org_id="org_123",
    role=role
)

```

## Go SDK

```go
resp, err := scalekitClient.Role().CreateOrganizationRole(ctx, "org_123", &rolesv1.CreateOrganizationRole{
	Name:        "org_admin",
	DisplayName: "Org Admin",
	Description: proto.String("Organization-scoped role"), // optional
	Extends:     proto.String("base_role_name"),        // optional
	Permissions: []string{"perm.read", "perm.write"},   // optional
})

```

## Java SDK

```java
CreateOrganizationRoleResponse res = scalekitClient.roles().createOrganizationRole(
    "org_123",
    CreateOrganizationRoleRequest.newBuilder()
        .setOrgId("org_123")
        .setRole(
            CreateOrganizationRole.newBuilder()
                .setName("org_admin")
                .setDisplayName("Org Admin")
                .setDescription("Organization-scoped role")
                .setExtends("base_role_name")          // optional
                .addPermissions("perm.read")           // optional
                .addPermissions("perm.write")          // optional
                .build()
        )
        .build()
);

```
