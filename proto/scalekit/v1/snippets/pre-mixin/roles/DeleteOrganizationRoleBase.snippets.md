---
operationId: RolesService_DeleteOrganizationRoleBase
---

```javascript
await scalekit.role.deleteOrganizationRoleBase("org_123", "senior_editor");
```

```python
scalekit_client.roles.delete_organization_role_base(
    org_id="org_123",
    role_name="senior_editor"
)
```

```go
err := scalekitClient.Role().DeleteOrganizationRoleBase(ctx, "org_123", "senior_editor")
if err != nil {
    // handle error
}
```

```java
scalekitClient.roles().deleteOrganizationRoleBase("org_123", "senior_editor");
```
