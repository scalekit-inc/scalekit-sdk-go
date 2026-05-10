---
operationId: ApiTokenService_InvalidateToken
---

```javascript
await scalekit.token.invalidateToken("apit_123456789012345");
```

```python
scalekit_client.token.invalidate_token(token="apit_123456789012345")
```

```go
err := scalekitClient.Token().InvalidateToken(ctx, "apit_123456789012345")
if err != nil {
    // handle error
}
```

```java
scalekitClient.tokens().invalidate("apit_123456789012345");
```
