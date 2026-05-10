---
operationId: ConnectionService_CreateConnection
---

```javascript
const connection = await scalekit.connection.createConnection(organizationId, connectionConfig);
```

```python
connection = scalekit_client.connection.create_connection(
  organization_id,
  connection_config
)
```

## Go SDK

```go
connection, err := scalekitClient.Connection.CreateConnection(
  ctx,
  organizationId,
  connectionConfig,
)
```

## Java SDK

```java
Connection connection = scalekitClient.connections().createConnection(organizationId, connectionConfig);
```
