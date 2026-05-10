---
operationId: DirectoryService_ListDirectories
---

```javascript
await scalekit.directory.listDirectories('<organization_id>');
```

```python
directories_list = scalekit_client.directory.list_directories(
	organization_id='<organization_id>'
)
```

## Go SDK

```go
directories,err := scalekitClient.Directory().ListDirectories(ctx, organizationId)
```

## Java SDK

```java
ListDirectoriesResponse response = scalekitClient.directories().listDirectories(organizationId);
```
