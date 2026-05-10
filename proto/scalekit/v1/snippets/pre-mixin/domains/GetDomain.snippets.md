---
operationId: DomainService_GetDomain
---

```js
// Fetch details of a specific domain
const response = await scalekit.domain.getDomain(organizationId, domainId);

// Domain object properties:
// - id: Domain identifier
// - domain: Domain name
// - organizationId: Owning organization
// - domainType: Domain configuration type
```

```py
# Fetch details of a specific domain
response = scalekit_client.domain.get_domain(organization_id="org_123", domain_id="dom_123")
```

```go
domain, err := scalekitClient.Domain().GetDomain(ctx, "dom_123", "org_123")
```

```java
Domain domain = scalekitClient.domains().getDomainById("org_123", "dom_123");
```
