---
operationId: RolesService_UpdateDefaultOrganizationRoles
---

```javascript
const res = await scalekit.role.updateDefaultOrganizationRoles("org_123", {
  defaultMemberRole: "org_member"
});
```

```python
from scalekit.v1.roles.roles_pb2 import UpdateDefaultOrganizationRolesRequest

res = scalekit_client.roles.update_default_organization_roles(
    org_id="org_123",
    default_roles=UpdateDefaultOrganizationRolesRequest(
        default_member_role="org_member"
    )
)
```

## Go SDK

```go
resp, err := scalekitClient.Role().UpdateDefaultOrganizationRoles(ctx, "org_123", "org_member")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdateDefaultOrganizationRolesResponse res = scalekitClient.roles().updateDefaultOrganizationRoles(
    "org_123",
    UpdateDefaultOrganizationRolesRequest.newBuilder()
        .setDefaultMemberRole("org_member")
        .build()
);
```
