---
operationId: DirectoryService_ListDirectoryGroups
---

```javascript
const { groups } = await scalekit.directory.listDirectoryGroups(
  '<organization_id>',
  '<directory_id>'
);
```

```python
directory_groups = scalekit_client.directory.list_directory_groups(
  directory_id='<directory_id>', organization_id='<organization_id>'
)
```

## Go SDK

```go
options := &ListDirectoryGroupsOptions{
		PageSize: 10,
		PageToken:"",
	}

directoryGroups, err := scalekitClient.Directory().ListDirectoryGroups(ctx, organizationId, directoryId, options)
```

## Java SDK

```java
var options = ListDirectoryResourceOptions.builder()
  .pageSize(10)
  .pageToken("")
  .includeDetail(true)
  .build();

ListDirectoryGroupsResponse groupsResponse = scalekitClient
  .directories()
  .listDirectoryGroups(directory.getId(), organizationId, options);
```
