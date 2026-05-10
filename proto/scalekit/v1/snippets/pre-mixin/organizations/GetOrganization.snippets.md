---
operationId: OrganizationService_GetOrganization
---

```javascript
const scalekit = new ScalekitClient(
  <SCALEKIT_ENVIRONMENT_URL>,
  <SCALEKIT_CLIENT_ID>,
  <SCALEKIT_CLIENT_SECRET>
);

const organization = await scalekit.organization.getOrganization(organization_id);
```

```python
scalekit_client = ScalekitClient(
  <SCALEKIT_ENVIRONMENT_URL>,
  <SCALEKIT_CLIENT_ID>,
  <SCALEKIT_CLIENT_SECRET>
)

organization = scalekit_client.organization.get_organization(
  organization_id
)
```

## Go SDK

```go
scalekitClient := scalekit.NewScalekitClient(
  <SCALEKIT_ENVIRONMENT_URL>,
  <SCALEKIT_CLIENT_ID>,
  <SCALEKIT_CLIENT_SECRET>
)

organization, err := scalekitClient.Organization.GetOrganization(
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

Organization organization = scalekitClient.organizations().getById(organizationId);
```
