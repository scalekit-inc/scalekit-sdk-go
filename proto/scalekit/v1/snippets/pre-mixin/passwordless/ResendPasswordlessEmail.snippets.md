---
operationId: PasswordlessService_ResendPasswordlessEmail
---

```js
const { authRequestId } = sendResponse;
const resendResponse = await scalekit.passwordless.resendPasswordlessEmail(
	authRequestId
);
```

```python
resend_response = scalekit_client.passwordless.resend_passwordless_email(
    auth_request_id=auth_request_id,
)

# New auth request ID from resend
new_auth_request_id = resend_response[0].auth_request_id
```

```go
resendResponse, err := scalekitClient.Passwordless().ResendPasswordlessEmail(
    ctx,
    authRequestId,
)

if err != nil {
    // Handle error
    return
}
```

```java
SendPasswordlessResponse resendResponse = passwordlessClient.resendPasswordlessEmail(authRequestId);
```
