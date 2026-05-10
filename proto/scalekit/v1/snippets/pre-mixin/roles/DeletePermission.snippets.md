---
operationId: RolesService_DeletePermission
---

```javascript
await scalekit.permission.deletePermission("read:users");
```

```python
scalekit_client.permissions.delete_permission(
    permission_name="read:users"
)
```

## Go SDK

```go
err := scalekitClient.Permission().DeletePermission(ctx, "read:users")
if err != nil { /* handle err */ }
```

## Java SDK

```java
scalekitClient.permissions().deletePermission("read:users");
```
