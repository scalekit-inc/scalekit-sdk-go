---
operationId: ConnectionService_DeleteConnection
---

```javascript
await scalekit.connection.deleteConnection(organizationId, connectionId);
```

```python
scalekit_client.connection.delete_connection(
  organization_id,
  connection_id
)
```

## Go SDK

```go
err := scalekitClient.Connection.DeleteConnection(
  ctx,
  organizationId,
  connectionId,
)
```

## Java SDK

```java
scalekitClient.connections().deleteConnection(connectionId, organizationId);
```
