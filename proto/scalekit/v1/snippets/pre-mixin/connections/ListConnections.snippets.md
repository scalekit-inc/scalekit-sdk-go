---
operationId: ConnectionService_ListConnections
---

```javascript
// List connections by organization id
const connections = await scalekit.connection.listConnections(organizationId);

// List connections by domain
const connections = await scalekit.connection.listConnectionsByDomain(domain);
```

```python
# List connections by organization id
connections = scalekit_client.connection.list_connections(
  organization_id
)

# List connections by domain
response = scalekit_client.connection.list_connections_by_domain(domain="example.com")
```

## Go SDK

```go
// List connections by organization id
connections, err := scalekitClient.Connection().ListConnections(
  ctx,
  organizationId
)

// List connections by domain
connections, err := scalekitClient.Connection().ListConnectionsByDomain(ctx, 
  domain)
```

## Java SDK

```java
// List connections by organization id
ListConnectionsResponse response = scalekitClient.connections(
  ).listConnections(organizationId);

// List connections by domain
ListConnectionsResponse response = scalekitClient.connections(
  ).listConnectionsByDomain("your-domain.com");
```
