---
operationId: UserService_ListUsers
---

```javascript
const response = await scalekit.user.listUsers(
  { pageSize: 100 });
```

```python
# pass empty org to fetch all users in environment
resp,_ = scalekit_client.users.list_users(organization_id="", page_size=100)
```

## Go SDK

```go
all, err := scalekitClient.User().ListUsers(ctx, &scalekit.ListUsersOptions{PageSize: 100})
```

## Java SDK

```java
ListUsersRequest lur = ListUsersRequest.
  newBuilder().setPageSize(100).build();
ListUsersResponse allUsers = users.listUsers(lur);
```
