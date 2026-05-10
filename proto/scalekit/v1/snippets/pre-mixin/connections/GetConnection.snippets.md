---
operationId: ConnectionService_GetConnection
---

```javascript
const connection = await scalekit.connection.getConnection(
  organizationId,
  connectionId
);
```

```python
connection = scalekit_client.connection.get_connection(
  organization_id,
  connection_id,
)
```

## Go SDK

```go
connection, err := scalekitClient.Connection.GetConnection(
  ctx,
  organizationId,
  connectionId,
)
```

## Java SDK

```java
Connection connection = scalekitClient.connections().getConnectionById(connectionId, organizationId);
```
