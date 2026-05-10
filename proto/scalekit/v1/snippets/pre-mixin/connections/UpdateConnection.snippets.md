---
operationId: ConnectionService_UpdateConnection
---

```javascript
const connection = await scalekit.connection.updateConnection(organizationId, connectionId, connectionConfig);
```

```python
connection = scalekit_client.connection.update_connection(
  organization_id,
  connection_id,
  connection_config
)
```

## Go SDK

```go
connection, err := scalekitClient.Connection.UpdateConnection(
  ctx,
  organizationId,
  connectionId,
  connectionConfig,
)
```

## Java SDK

```java
Connection connection = scalekitClient.connections().updateConnection(connectionId, organizationId, connectionConfig);
```
