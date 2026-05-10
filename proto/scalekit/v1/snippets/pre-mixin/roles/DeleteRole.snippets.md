---
operationId: RolesService_DeleteRole
---

```javascript
// Basic delete
await scalekit.role.deleteRole("admin");

// With reassignment
await scalekit.role.deleteRole("admin", "member");
```

```python
# Basic delete
scalekit_client.roles.delete_role(role_name="admin")

# With reassignment
scalekit_client.roles.delete_role(
    role_name="admin",
    reassign_role_name="member"
)
```

## Go SDK

```go
// Basic delete
err := scalekitClient.Role().DeleteRole(ctx, "admin")
if err != nil { /* handle err */ }

// With reassignment
err = scalekitClient.Role().DeleteRole(ctx, "admin", "member")
```

## Java SDK

```java
// Basic delete
scalekitClient.roles().deleteRole("admin");

// With reassignment
scalekitClient.roles().deleteRole("admin", "member");
```
