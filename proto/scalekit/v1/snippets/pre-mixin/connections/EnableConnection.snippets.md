---
operationId: ConnectionService_EnableConnection
---

```javascript
await scalekit.connection.enableConnection(organizationId, connectionId);
```

```python
scalekit_client.connection.enable_connection(
  organization_id,
  connection_id,
)
```

## Go SDK

```go
err := scalekitClient.Connection.EnableConnection(
  ctx,
  organizationId,
  connectionId,
)
```

## Java SDK

```java
ToggleConnectionResponse response = scalekitClient.connections().enableConnection(connectionId, organizationId);
```
