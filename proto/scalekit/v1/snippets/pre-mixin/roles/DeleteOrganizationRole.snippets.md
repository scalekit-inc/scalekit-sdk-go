---
operationId: RolesService_DeleteOrganizationRole
---

```javascript
// Basic delete
await scalekit.role.deleteOrganizationRole("org_123", "org_role_admin");

// With reassignment
await scalekit.role.deleteOrganizationRole("org_123", "org_role_admin", "org_role_member");
```

```python
# Basic delete
scalekit_client.roles.delete_organization_role(
    org_id="org_123",
    role_name="org_role_admin"
)

# With reassignment
scalekit_client.roles.delete_organization_role(
    org_id="org_123",
    role_name="org_role_admin",
    reassign_role_name="org_role_member"
)

```

## Go SDK

```go
// Basic delete
err := scalekitClient.Role().DeleteOrganizationRole(ctx, "org_123", "org_role_admin")
if err != nil { /* handle err */ }

// With reassignment
err = scalekitClient.Role().DeleteOrganizationRole(ctx, "org_123", "org_role_admin", "org_role_member")

```

## Java SDK

```java
// Basic delete
scalekitClient.roles().deleteOrganizationRole("org_123", "org_role_admin");

// With reassignment
scalekitClient.roles().deleteOrganizationRole("org_123", "org_role_admin", "org_role_member");

```
