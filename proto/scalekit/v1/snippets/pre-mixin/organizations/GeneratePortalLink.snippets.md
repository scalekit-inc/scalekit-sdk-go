---
operationId: OrganizationService_GeneratePortalLink
---

```javascript
const link = await scalekit.organization.generatePortalLink(organizationId);
```

```python
link = scalekit_client.organization.generate_portal_link(
  organization_id
)
```

## Go SDK

```go
link, err := scalekitClient.Organization.GeneratePortalLink(
  ctx,
  organizationId
)
```

## Java SDK

```java
Link portalLink = client
  .organizations()
  .generatePortalLink(organizationId, Arrays.asList(Feature.sso, Feature.dir_sync));
```
