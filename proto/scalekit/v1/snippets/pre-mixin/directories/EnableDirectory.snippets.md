---
operationId: DirectoryService_EnableDirectory
---

```javascript
await scalekit.directory.enableDirectory('<organization_id>', '<directory_id>');
```

```python
directory_response = scalekit_client.directory.enable_directory(
  directory_id='<directory_id>', organization_id='<organization_id>'
)
```

## Go SDK

```go
enable,err := scalekitClient.Directory().EnableDirectory(ctx, organizationId, directoryId)
```

## Java SDK

```java
ToggleDirectoryResponse enableResponse = client
  .directories()
  .enableDirectory(directoryId, organizationId);
```
