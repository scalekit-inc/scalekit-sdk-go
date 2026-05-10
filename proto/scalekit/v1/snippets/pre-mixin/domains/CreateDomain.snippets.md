---
operationId: DomainService_CreateDomain
---

```js
// Add a new domain to an organization
const response = await scalekit.createDomain("org-123", "example.com", {
	// Domain type: controls user authentication and email validation
	domainType: "ORGANIZATION_DOMAIN",
});
```

```python
# Add a new domain to an organization
response = scalekit_client.domain.create_domain(organization_id="org-123",
			domain_name="example.com",
 			domain_type="ORGANIZATION_DOMAIN")
```

```go
domain, err := scalekitClient.Domain().CreateDomain(ctx, "org_id", "example.com", &scalekit.CreateDomainOptions{
		DomainType: "ORGANIZATION_DOMAIN",
	})
```

```java
CreateDomainRequest request = CreateDomainRequest.newBuilder()
	.setOrganizationId(organization.getId())
	.setDomain(CreateDomain.newBuilder()
		.setDomain("example.com")
		.setDomainType("ORGANIZATION_DOMAIN")
		.build())
	.build();
	
Domain domain = scalekitClient.domains().createDomain(request);
```
