---
operationId: RolesService_ListOrganizationRoles
---

```javascript
const res = await scalekit.role.listOrganizationRoles("org_123");
```

```python
res = scalekit_client.roles.list_organization_roles(org_id="org_123")
```

## Go SDK

```go
resp, err := scalekitClient.Role().ListOrganizationRoles(ctx, "org_123")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
ListOrganizationRolesResponse res = scalekitClient.roles().listOrganizationRoles("org_123");
```
