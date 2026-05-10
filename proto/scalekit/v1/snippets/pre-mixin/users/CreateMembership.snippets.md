---
operationId: UserService_CreateMembership
---

```javascript
import { ScalekitClient } from "@scalekit-sdk/node";
const scalekit = new ScalekitClient(
	process.env.SCALEKIT_ENV_URL,
	process.env.SCALEKIT_CLIENT_ID,
	process.env.SCALEKIT_CLIENT_SECRET
);
await scalekit.user.createMembership("org_123", "usr_123", {
	roles: ["admin"],
	metadata: {
		department: "engineering",
		location: "nyc-office",
	},
});
```

```python
from scalekit.v1.users.users_pb2 import CreateMembership
from scalekit.v1.commons.commons_pb2 import Role

membership = CreateMembership(
    roles=[Role(name="admin")],
    metadata={"department": "engineering", "location": "nyc-office"},
)
resp = scalekit_client.users.create_membership(
    organization_id="org_123",
    user_id="usr_123",
    membership=membership,
)
```

## Go SDK

```go
func main() {
    scalekitClient := scalekit.NewScalekitClient(
        os.Getenv("SCALEKIT_ENV_URL"),
        os.Getenv("SCALEKIT_CLIENT_ID"),
        os.Getenv("SCALEKIT_CLIENT_SECRET"),
    )
    membership := &usersv1.CreateMembership{
        Roles: []*usersv1.Role{{Name: "admin"}},
        Metadata: map[string]string{
            "department": "engineering",
            "location":   "nyc-office",
        },
    }
    resp, 
      err := scalekitClient.User().CreateMembership(
        context.Background(), "org_123", 
          "usr_123", membership, false)
    if err != nil {
        panic(err)
    }
}
```

## Java SDK

```java
import com.scalekit.ScalekitClient;
import com.scalekit.api.UserClient;
import com.scalekit.grpc.scalekit.v1.users.*;
ScalekitClient scalekitClient = new ScalekitClient(
    System.getenv("SCALEKIT_ENV_URL"),
    System.getenv("SCALEKIT_CLIENT_ID"),
    System.getenv("SCALEKIT_CLIENT_SECRET")
);
UserClient users = scalekitClient.users();
CreateMembershipRequest membershipReq = CreateMemb
  ershipRequest.newBuilder()
        .setMembership(
          CreateMembership.newBuilder()
                .addRoles(Role.newBuilder(
                  ).setName("admin").build())
                .putMetadata("department", "engineering")
                .putMetadata("location", "nyc-office")
                .build())
        .build();
CreateMembershipResponse res = users.
  createMembership("org_123", "usr_123", 
    membershipReq);
```
