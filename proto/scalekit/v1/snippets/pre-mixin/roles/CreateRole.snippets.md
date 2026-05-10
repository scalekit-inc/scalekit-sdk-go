---
operationId: RolesService_CreateRole
---

```javascript
await scalekit.role.createRole({
  name: "admin",
  displayName: "Admin",
  description: "Environment-level role",
  extends: "base_role",                // optional
  permissions: ["read:users"]          // optional
});

```

```python
from scalekit.v1.roles.roles_pb2 import CreateRole

role = CreateRole(
    name="admin",
    display_name="Admin",
    description="Environment-level role",
    extends="base_role",                  # optional
    permissions=["read:users"]           # optional
)

scalekit_client.roles.create_role(role=role)

```

## Go SDK

```go
resp, err := scalekitClient.Role().CreateRole(ctx, &rolesv1.CreateRole{
	Name:        "admin",
	DisplayName: "Admin",
	Description: "Environment-level role",
	Extends:     proto.String("base_role"),      // optional
	Permissions: []string{"read:users"},        // optional
})


```

## Java SDK

```java
CreateRoleResponse res = scalekitClient.roles().createRole(
    CreateRoleRequest.newBuilder()
        .setRole(
            CreateRole.newBuilder()
                .setName("admin")
                .setDisplayName("Admin")
                .setDescription("Environment-level role")
                // .setExtends("base_role")         // optional
                // .addPermissions("read:users")    // optional
                .build()
        )
        .build()
);

```
