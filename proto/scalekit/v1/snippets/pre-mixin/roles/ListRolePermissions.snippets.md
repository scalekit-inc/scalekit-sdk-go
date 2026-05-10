---
operationId: RolesService_ListRolePermissions
---

```javascript
const res = await scalekit.permission.listRolePermissions("admin");
```

```python
res = scalekit_client.permissions.list_role_permissions(role_name="admin")
```

## Go SDK

```go
resp, err := scalekitClient.Permission().ListRolePermissions(ctx, "admin")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
ListRolePermissionsResponse res = scalekitClient.permissions().listRolePermissions("admin");
```
