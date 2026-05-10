---
operationId: RolesService_GetRole
---

```javascript
const res = await scalekit.role.getRole("admin");
```

```python
res = scalekit_client.roles.get_role(role_name="admin")
```

## Go SDK

```go
resp, err := scalekitClient.Role().GetRole(ctx, "admin")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
GetRoleResponse res = scalekitClient.roles().getRole("admin");
```
