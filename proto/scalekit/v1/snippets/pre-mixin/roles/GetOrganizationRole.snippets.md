---
operationId: RolesService_GetOrganizationRole
---

```javascript
const res = await scalekit.role.getOrganizationRole("org_123", "org_admin");
```

```python
res = scalekit_client.roles.get_organization_role(
    org_id="org_123",
    role_name="org_admin"
)
```

## Go SDK

```go
resp, err := scalekitClient.Role().GetOrganizationRole(ctx, "org_123", "org_admin")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
GetOrganizationRoleResponse res = scalekitClient.roles().getOrganizationRole("org_123", "org_admin");
```
