---
operationId: UserService_CreateUserAndMembership
---

```javascript
const {
   user } = await scalekit.user.
    createUserAndMembership("org_123", {
	email: "user@example.com",
	externalId: "ext_12345a67b89c",
	metadata: { department: "engineering", 
	  location: "nyc-office" },
	userProfile: {
		firstName: "John",
		lastName: "Doe",
	},
});
```

```python
# Create user with membership 
user = CreateUser(
    email="john.doe@example.com",
    external_id="ext_john_123",  # Optional
    user_profile={
        "first_name": "John",
        "last_name": "Doe",
        "name": "John Doe",
        "locale": "en-US",
        "phone_number": "+14155552671"
    },
    membership={
        "roles": [{"name": "member"}]  
    }
)

# Create user and membership in organization
response = scalekit_client.users.create_user_and_membership(
    organization_id="your_org_id",
    user=user,
    send_invitation_email=True  # Set to False if you don't want to send email
)

    user_id = response[0].user.id
```

## Go SDK

```go
newUser := &usersv1.CreateUser{
    Email:      "user@example.com",
    ExternalId: "ext_12345a67b89c",
    Metadata: map[string]string{
        "department": "engineering",
        "location":   "nyc-office",
    },
    UserProfile: &usersv1.CreateUserProfile{
        FirstName: "John",
        LastName:  "Doe",
    },
}
cuResp, 
  err := scalekitClient.User().CreateUserAndMembership(ctx, "org_123", newUser, false)
if err != nil { /* handle error */ }
```

## Java SDK

```java
CreateUser createUser = CreateUser.newBuilder()
        .setEmail("user@example.com")
        .setExternalId("ext_12345a67b89c")
        .putMetadata("department", "engineering")
        .putMetadata("location", "nyc-office")
        .setUserProfile(
          CreateUserProfile.newBuilder()
                .setFirstName("John")
                .setLastName("Doe")
                .build())
        .build();
CreateUserAndMembershipRequest cuReq = CreateUserA
  ndMembershipRequest.newBuilder()
        .setUser(createUser)
        .build();
CreateUserAndMembershipResponse cuResp = users.
  createUserAndMembership("org_123", cuReq);
System.out.println(cuResp.getUser().getId());
```
