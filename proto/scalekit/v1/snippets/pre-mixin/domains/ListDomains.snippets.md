---
operationId: DomainService_ListDomains
---

```js
// List all domains in an organization
const response = await scalekit.domain.listDomains(organizationId, {
	domainType: "ORGANIZATION_DOMAIN"
});

// Domain object contains:
// - id: Domain identifier
// - domain: Domain name
// - organizationId: Owning organization
// - domainType: Configuration type
```

```py
# List all domains in an organization
response = scalekit_client.domain.list_domains(
            organization_id="org_123",
            domain_type="ORGANIZATION_DOMAIN"
        )
# - organization_id: Owning organization
# - domain_type: domain type
```

```go
domains, err := scalekitClient.Domain().ListDomains(ctx, "org_id", &scalekit.ListDomainOptions{
DomainType: "ORGANIZATION_DOMAIN",
})
```

```java
List<Domain> domains = scalekitClient.domains().listDomainsByOrganizationId("org_id", "ORGANIZATION_DOMAIN");
```
