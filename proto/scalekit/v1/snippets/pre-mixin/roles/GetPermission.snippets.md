---
operationId: RolesService_GetPermission
---

```javascript
const res = await scalekit.permission.getPermission("read:users");
```

```python
res = scalekit_client.permissions.get_permission(
    permission_name="read:users"
)
```

## Go SDK

```go
resp, err := scalekitClient.Permission().GetPermission(ctx, "read:users")
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
GetPermissionResponse res = scalekitClient.permissions().getPermission("read:users");
```
