---
operationId: DirectoryService_ListDirectoryUsers
---

```javascript
const { users } = await scalekit.directory.listDirectoryUsers(
  '<organization_id>',
  '<directory_id>'
);
```

```python
directory_users = scalekit_client.directory.list_directory_users(
  directory_id='<directory_id>', organization_id='<organization_id>'
)
```

## Go SDK

```go
options := &ListDirectoryUsersOptions{
		PageSize: 10,
		PageToken: "",
	}
directoryUsers,err := scalekitClient.Directory().ListDirectoryUsers(ctx, organizationId, directoryId, options)
```

## Java SDK

```java
var options = ListDirectoryResourceOptions.builder()
  .pageSize(10)
  .pageToken("")
  .includeDetail(true)
  .build();

ListDirectoryUsersResponse usersResponse = scalekitClient
  .directories()
  .listDirectoryUsers(directory.getId(), organizationId, options);
```
