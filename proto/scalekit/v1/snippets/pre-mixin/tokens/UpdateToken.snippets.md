---
operationId: ApiTokenService_UpdateToken
---

```javascript
// Merge custom claims (set value to "" to remove a claim)
const { tokenInfo } = await scalekit.token.updateToken("apit_123456789012345", {
  customClaims: { env: "staging", old_claim: "" },
  description: "Updated CI Token",
});
```

```python
# Merge custom claims (set value to "" to remove a claim)
resp = scalekit_client.token.update_token(
    token="apit_123456789012345",
    custom_claims={"env": "staging", "old_claim": ""},
    description="Updated CI Token",
)
```

```go
// Merge custom claims (set value to "" to remove a claim)
// Description: nil leaves it unchanged, pointer to "" clears it
desc := "Updated CI Token"
resp, err := scalekitClient.Token().UpdateToken(ctx, "apit_123456789012345", scalekit.UpdateTokenOptions{
    CustomClaims: map[string]string{"env": "staging", "old_claim": ""},
    Description:  &desc,
})
if err != nil {
    // handle error
}
```

```java
// Merge custom claims (set value to "" to remove a claim)
Map<String, String> claims = Map.of("env", "staging", "old_claim", "");
UpdateTokenResponse resp = scalekitClient.tokens().update(
    "apit_123456789012345", claims, "Updated CI Token"
);
```
