---
operationId: PasswordlessService_SendPasswordlessEmail
---

```js
const response = await scalekit.passwordless.sendPasswordlessEmail(
	"john.doe@example.com",
	{
		template: "SIGNIN",
		expiresIn: 100,
		magiclinkAuthUri: "https://www.google.com",
		templateVariables: {
			employeeID: "EMP523",
			teamName: "Alpha Team",
			supportEmail: "support@yourcompany.com",
		},
	}
);
```

```python
response = scalekit_client.passwordless.send_passwordless_email(
    email="john.doe@example.com",
    template="SIGNIN",  # or "SIGNUP", "UNSPECIFIED"
    expires_in=100,
    magiclink_auth_uri="https://www.google.com",
    template_variables={
        "employeeID": "EMP523",
        "teamName": "Alpha Team",
        "supportEmail": "support@yourcompany.com",
    },
)

# Extract auth request ID from response
auth_request_id = response[0].auth_request_id
```

```go
templateType := scalekit.TemplateTypeSignin
response, err := scalekitClient.Passwordless().SendPasswordlessEmail(
    ctx,
    "john.doe@example.com",
    &scalekit.SendPasswordlessOptions{
        Template:         &templateType,
        ExpiresIn:        100,
        MagiclinkAuthUri: "https://www.google.com",
        TemplateVariables: map[string]string{
            "employeeID":    "EMP523",
            "teamName":      "Alpha Team",
            "supportEmail":  "support@yourcompany.com",
        },
    },
)

if err != nil {
    // Handle error
    return
}

authRequestId := response.AuthRequestId
```

```java
TemplateType templateType = TemplateType.SIGNIN;
Map<String, String> templateVariables = new HashMap<>();
templateVariables.put("employeeID", "EMP523");
templateVariables.put("teamName", "Alpha Team");
templateVariables.put("supportEmail", "support@yourcompany.com");

SendPasswordlessOptions options = new SendPasswordlessOptions();
options.setTemplate(templateType);
options.setExpiresIn(100);
options.setMagiclinkAuthUri("https://www.example.com");
options.setTemplateVariables(templateVariables);

SendPasswordlessResponse response = passwordlessClient.sendPasswordlessEmail(
    "john.doe@example.com",
    options
);

String authRequestId = response.getAuthRequestId();
```
