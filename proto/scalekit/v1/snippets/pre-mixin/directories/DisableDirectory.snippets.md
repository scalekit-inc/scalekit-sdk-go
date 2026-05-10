---
operationId: DirectoryService_DisableDirectory
---

```javascript
await scalekit.directory.disableDirectory(
  '<organization_id>',
  '<directory_id>'
);
```

```python
directory_response = scalekit_client.directory.disable_directory(
  directory_id='<directory_id>', organization_id='<organization_id>'
)
```

## Go SDK

```go
disable,err := scalekitClient.Directory().DisableDirectory(ctx, organizationId, directoryId)
```

## Java SDK

```java
ToggleDirectoryResponse disableResponse = scalekitClient
  .directories()
  .disableDirectory(directoryId, organizationId);
```
