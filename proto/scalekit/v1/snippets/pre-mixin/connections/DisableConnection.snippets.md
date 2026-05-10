---
operationId: ConnectionService_DisableConnection
---

```javascript
await scalekit.connection.disableConnection(organizationId, connectionId);
```

```python
scalekit_client.connection.disable_connection(
  organization_id,
  connection_id
)
```

## Go SDK

```go
err := scalekitClient.Connection.DisableConnection(
  ctx,
  organizationId,
  connectionId,
)
```

## Java SDK

```java
ToggleConnectionResponse response = scalekitClient.connections().disableConnection(connectionId, organizationId);
```
