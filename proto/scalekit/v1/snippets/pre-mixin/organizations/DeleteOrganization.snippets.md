---
operationId: OrganizationService_DeleteOrganization
---

```javascript
await scalekit.organization.deleteOrganization(organizationId);
```

```python
scalekit_client.organization.delete_organization(organization_id)
```

## Go SDK

```go
err := scalekitClient.Organization.DeleteOrganization(
  ctx,
  organizationId
)
```

## Java SDK

```java
ScalekitClient scalekitClient = new ScalekitClient(
  "<SCALEKIT_ENVIRONMENT_URL>",
  "<SCALEKIT_CLIENT_ID>",
  "<SCALEKIT_CLIENT_SECRET>"
);

scalekitClient.organizations().deleteById(organizationId);
```
