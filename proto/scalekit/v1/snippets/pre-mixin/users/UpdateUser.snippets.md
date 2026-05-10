---
operationId: UserService_UpdateUser
---

```javascript
await scalekit.user.updateUser("usr_123", {
	userProfile: {
		firstName: "John",
		lastName: "Smith",
	},
	metadata: {
		department: "sales",
	},
});
```

```python
import os
from scalekit import ScalekitClient
from scalekit.v1.users.users_pb2 import UpdateUser
from scalekit.v1.commons.commons_pb2 import UserProfile
scalekit_client = ScalekitClient(
    env_url=os.getenv("SCALEKIT_ENV_URL"),
    client_id=os.getenv("SCALEKIT_CLIENT_ID"),
    client_secret=os.getenv("SCALEKIT_CLIENT_SECRET"),
)
update_user = UpdateUser(
    user_profile=UserProfile(
        first_name="John",
        last_name="Smith"
    ),
    metadata={"department": "sales"}
)
scalekit_client.users.update_user(organization_id="org_123", 
  user=update_user)
```

## Go SDK

```go
upd := &usersv1.UpdateUser{
    UserProfile: &usersv1.UpdateUserProfile{
        FirstName: "John",
        LastName:  "Smith",
    },
    Metadata: map[string]string{
        "department": "sales",
    },
}
scalekitClient.User().UpdateUser(ctx, "usr_123", upd)
```

## Java SDK

```java
UpdateUser upd = UpdateUser.newBuilder()
        .setUserProfile(
          UpdateUserProfile.newBuilder()
                .setFirstName("John")
                .setLastName("Smith")
                .build())
        .putMetadata("department", "sales")
        .build();
UpdateUserRequest updReq = UpdateUserRequest.
  newBuilder().setUser(upd).build();
users.updateUser("usr_123", updReq);
```
