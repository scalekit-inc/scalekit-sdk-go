---
operationId: SessionService_GetUserSessions
---

```javascript
// Basic usage
const res = await scalekit.session.getUserSessions("user_123");

// With pagination and filtering
const res = await scalekit.session.getUserSessions("user_123", {
  pageSize: 10,
  pageToken: "next_page_token",
  filter: {
    status: ["ACTIVE"],
    startTime: new Date("2024-01-01"),
    endTime: new Date("2024-12-31")
  }
});
```

```python
# Basic usage
res = scalekit_client.sessions.get_user_sessions(user_id="user_123")

# With pagination and filtering
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime

start_time = Timestamp()
start_time.FromDatetime(datetime(2024, 1, 1))
end_time = Timestamp()
end_time.FromDatetime(datetime(2024, 12, 31))

filter_obj = scalekit_client.sessions.create_session_filter(
    status=["ACTIVE"], start_time=start_time, end_time=end_time
)
res = scalekit_client.sessions.get_user_sessions(
    user_id="user_123", page_size=10, page_token="next_page_token", filter=filter_obj
)
```

## Go SDK

```go
// Basic usage
resp, err := scalekitClient.Session().GetUserSessions(ctx, "user_123", 0, "", nil)
if err != nil { /* handle err */ }

// With pagination and filtering
// import "time", sessionsv1 "...", "google.golang.org/protobuf/types/known/timestamppb"
startTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
endTime, _ := time.Parse(time.RFC3339, "2024-12-31T23:59:59Z")
filter := &sessionsv1.UserSessionFilter{
    Status:    []string{"ACTIVE"},
    StartTime: timestamppb.New(startTime),
    EndTime:   timestamppb.New(endTime),
}
resp, err := scalekitClient.Session().GetUserSessions(ctx, "user_123", 10, "next_page_token", filter)
if err != nil { /* handle err */ }
```

## Java SDK

```java
// Basic usage
UserSessionDetails res = scalekitClient.sessions().getUserSessions("user_123", null, null, null);

// With pagination and filtering
// import UserSessionFilter, Timestamp, Instant
UserSessionFilter filter = UserSessionFilter.newBuilder()
    .addStatus("ACTIVE")
    .setStartTime(Timestamp.newBuilder().setSeconds(Instant.parse("2024-01-01T00:00:00Z").getEpochSecond()).build())
    .setEndTime(Timestamp.newBuilder().setSeconds(Instant.parse("2024-12-31T23:59:59Z").getEpochSecond()).build())
    .build();
UserSessionDetails res = scalekitClient.sessions().getUserSessions("user_123", 10, "next_page_token", filter);
```
