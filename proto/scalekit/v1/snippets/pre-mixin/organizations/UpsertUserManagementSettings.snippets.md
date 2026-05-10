---
operationId: OrganizationsService_UpsertUserManagementSettings
---

```javascript
const response = await scalekit.organization.upsertUserManagementSettings(
  "org_123",
  {
    maxAllowedUsers: 100,
  }
);
// Settings updated: { maxAllowedUsers: 100 }
```

```python
response = scalekit_client.organization.upsert_user_management_settings(
    organization_id="org_123",
    max_allowed_users=100
)
# Settings updated: max_allowed_users = 100
```

```go
resp, err := scalekitClient.Organization().UpsertUserManagementSettings(ctx, "org_123", scalekit.UpsertUserManagementSettingsOptions{
    MaxAllowedUsers: 100,
})
if err != nil {
    // handle error
}
// Settings updated: maxAllowedUsers = 100
```

```java
UpsertUserManagementSettingsResponse resp = scalekitClient.organizations().upsertUserManagementSettings(
    "org_123",
    100  // maxAllowedUsers
);
// Settings updated: maxAllowedUsers = 100
```
