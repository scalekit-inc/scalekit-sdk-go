---
operationId: RolesService_GetRoleUsersCount
---

```javascript
const res = await scalekit.role.getRoleUsersCount("admin");
```

```python
res = scalekit_client.roles.get_role_users_count(role_name="admin")
```

## Go SDK

```go
resp, err := scalekitClient.Role().GetRoleUsersCount(ctx, "admin")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
GetRoleUsersCountResponse res = scalekitClient.roles().getRoleUsersCount("admin");
```
