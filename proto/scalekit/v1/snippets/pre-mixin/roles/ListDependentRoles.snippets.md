---
operationId: RolesService_ListDependentRoles
---

```javascript
const res = await scalekit.role.listDependentRoles("admin");
```

```python
res = scalekit_client.roles.list_dependent_roles(role_name="admin")
```

## Go SDK

```go
resp, err := scalekitClient.Role().ListDependentRoles(ctx, "admin")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
ListDependentRolesResponse res = scalekitClient.roles().listDependentRoles("admin");
```
