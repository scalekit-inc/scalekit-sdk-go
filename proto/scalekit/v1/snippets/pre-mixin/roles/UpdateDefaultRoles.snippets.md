---
operationId: RolesService_UpdateDefaultRoles
---

```javascript
const res = await scalekit.role.updateDefaultRoles({
  defaultMemberRole: "member"
});
```

```python
from scalekit.v1.roles.roles_pb2 import UpdateDefaultRolesRequest

res = scalekit_client.roles.update_default_roles(
    default_roles=UpdateDefaultRolesRequest(
        default_member_role="member"
    )
)
```

## Go SDK

```go
resp, err := scalekitClient.Role().UpdateDefaultRoles(ctx, "member")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
UpdateDefaultRolesResponse res = scalekitClient.roles().updateDefaultRoles(
    UpdateDefaultRolesRequest.newBuilder()
        .setDefaultMemberRole("member")
        .build()
);
```
