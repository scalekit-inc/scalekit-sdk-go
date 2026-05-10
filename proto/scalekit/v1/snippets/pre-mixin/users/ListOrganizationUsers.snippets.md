---
operationId: UserService_ListOrganizationUsers
---

```javascript
const response = await scalekit.user.
  listOrganizationUsers("org_123", {
	pageSize: 50,
});
console.log(response.users);
```

```python
resp, _ = scalekit_client.users.list_users(organization_id="org_123", page_size=50)
```

## Go SDK

```go
list, 
  err := scalekitClient.User().ListOrganizationUsers(ctx, "org_123", &scalekit.ListUsersOptions{PageSize: 50})
if err != nil { /* handle error */ }
fmt.Println(list.Users)
```

## Java SDK

```java
ListOrganizationUsersRequest listReq = ListOrganiz
  ationUsersRequest.newBuilder()
        .setPageSize(50)
        .build();
ListOrganizationUsersResponse list = users.
  listOrganizationUsers("org_123", listReq);
```
