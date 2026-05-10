---
operationId: OrganizationService_UpdateOrganization
---

```javascript
const organization = await scalekit.organization.updateOrganization(organization_id, {
  displayName: 'displayName',
  externalId: 'externalId',
});
```

```python
organization = scalekit_client.organization.update_organization(organization_id, {
  display_name: "display_name",
  external_id: "external_id"
})
```

## Go SDK

```go
organization, err := scalekitClient.Organization.UpdateOrganization(
  ctx,
  organizationId,
  &scalekit.UpdateOrganization{
    DisplayName: "displayName",
    ExternalId: "externalId",
  },
)
```

## Java SDK

```java
UpdateOrganization updateOrganization = UpdateOrganization.newBuilder()
  .setDisplayName("Updated Organization Name")
  .build();

Organization updatedOrganizationById = scalekitClient.organizations().updateById(organizationId, updateOrganization);
```
