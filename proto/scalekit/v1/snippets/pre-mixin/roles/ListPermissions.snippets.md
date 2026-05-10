---
operationId: RolesService_ListPermissions
---

```javascript
const res = await scalekit.permission.listPermissions();
```

```python
res = scalekit_client.permissions.list_permissions()
```

## Go SDK

```go
resp, err := scalekitClient.Permission().ListPermissions(ctx)
if err != nil { /* handle err */ }
_ = resp
```

## Java SDK

```java
ListPermissionsResponse res = scalekitClient.permissions().listPermissions();
```
