---
operationId: DirectoryService_GetDirectory
---

```javascript
const { directory } = await scalekit.directory.getDirectory(
  organizationId,
  directoryId
);
```

```python
directory = scalekit_client.directory.get_directory(
  directory_id='<directory_id>', organization_id='<organization_id>'
)
print(f'Directory details: {directory}')
```

## Go SDK

```go
directory, err := scalekitClient.Directory().GetDirectory(ctx, organizationId, directoryId)
```

## Java SDK

```java
Directory directory = scalekitClient.directories().getDirectory(directoryId, organizationId);
```
