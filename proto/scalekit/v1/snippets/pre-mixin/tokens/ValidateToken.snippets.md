---
operationId: ApiTokenService_ValidateToken
---

```javascript
const { tokenInfo } = await scalekit.token.validateToken("apit_123456789012345");
const { organizationId, customClaims } = tokenInfo;
```

```python
resp = scalekit_client.token.validate_token(token="apit_123456789012345")
organization_id = resp.token_info.organization_id
custom_claims = resp.token_info.custom_claims
```

```go
resp, err := scalekitClient.Token().ValidateToken(ctx, "apit_123456789012345")
if err != nil {
    // handle error — token invalid, expired, or not found
}
organizationId := resp.TokenInfo.OrganizationId
customClaims := resp.TokenInfo.CustomClaims
```

```java
ValidateTokenResponse resp = scalekitClient.tokens().validate("apit_123456789012345");
String organizationId = resp.getTokenInfo().getOrganizationId();
Map<String, String> customClaims = resp.getTokenInfo().getCustomClaimsMap();
```
