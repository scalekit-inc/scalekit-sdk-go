---
operationId: ApiTokenService_CreateToken
---

```javascript
const { token, tokenId, tokenInfo } = await scalekit.token.createToken("org_123", {
  customClaims: { env: "prod", scope: "read" },
  description: "CI Deploy Token",
});
// token   → pass to the system that needs to authenticate
// tokenId → store on your side for future management (validate, invalidate, update)
```

```python
resp = scalekit_client.token.create_token(
    organization_id="org_123",
    custom_claims={"env": "prod", "scope": "read"},
    description="CI Deploy Token",
)
# resp.token      → pass to the system that needs to authenticate
# resp.token_id   → store on your side for future management (validate, invalidate, update)
# resp.token_info → token metadata
```

```go
resp, err := scalekitClient.Token().CreateToken(ctx, "org_123", scalekit.CreateTokenOptions{
    CustomClaims: map[string]string{"env": "prod", "scope": "read"},
    Description:  "CI Deploy Token",
})
if err != nil {
    // handle error
}
// resp.Token   → pass to the system that needs to authenticate
// resp.TokenId → store on your side for future management (validate, invalidate, update)
```

```java
Map<String, String> claims = Map.of("env", "prod", "scope", "read");
CreateTokenResponse resp = scalekitClient.tokens().create(
    "org_123", null, claims, null, "CI Deploy Token"
);
// resp.getToken()   → pass to the system that needs to authenticate
// resp.getTokenId() → store on your side for future management (validate, invalidate, update)
```
