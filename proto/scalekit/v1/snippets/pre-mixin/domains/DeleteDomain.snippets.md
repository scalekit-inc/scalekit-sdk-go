---
operationId: DomainService_DeleteDomain
---

```js
// Remove a domain from an organization
// Caution: Deletion is permanent and may affect user access
const response = await scalekit.domain.deleteDomain(organizationId, domainId);
```

```py
# Remove a domain from an organization
# Caution: Deletion is permanent and may affect user access
response = scalekit_client.domain.delete_domain(
    organization_id="org_123",
    domain_id="dom_123"
)
```

```go
err = scalekitClient.Domain().DeleteDomain(ctx, "dom_123", "org_123")
```

```java
scalekitClient.domains().deleteDomain(organization.getId(), domain.getId());
```
