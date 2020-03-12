## Notifications

### GroupCreated

```json
{
    "id": 0, 
    "groupAvatar": File, 
    "creatorAvatar": File,
    "parentAvatar (optional)": File,
    "group": Group,
    "creator": User,
    "seen (optional)": false
}
```

### GroupCreationRequestStarted

```json
{
    "id": 0, 
    "initiatorAvatar": File,
    "parentAvatar": File,
    "groupAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupCreationRequestAccepted

```json
{
    "id": 0, 
    "acceptorAvatar": File,
    "group": Group,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupCreationRequestDenied

```json
{
    "id": 0, 
    "acceptorAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupCreationRequestWithdrawn

```json
{
    "id": 0, 
    "acceptorAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupJoinRequestStarted

```json
{
    "id": 0,
    "initiatorAvatar": File,
    "groupAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupJoinRequestAccepted

```json
{
    "id": 0,
    "acceptorAvatar": File,
    "groupAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupJoinRequestDenied

```json
{
    "id": 0,
    "acceptorAvatar": File,
    "groupAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### GroupJoinRequestWithdrawn

```json
{
    "id": 0,
    "acceptorAvatar": File,
    "groupAvatar": File,
    "request": GroupCreationRequest,
    "seen (optional)": false
}
```

### UserAddedToGroup

```json
{
    "id": 0,
    "group": Group,
    "groupAvatar": File,
    "added": User,
    "addedAvatar": File,
    "addedBy": User,
    "addedByAvatar": File,
    "seen (optional)": false
}
```

### UserLeftGroup

```json
{
    "id": 0,
    "group": Group,
    "groupAvatar": File,
    "user": User,
    "userAvatar": File,
    "seen (optional)": false
}
```

### UserRemovedFromGroup

```json
{
    "id": 0,
    "group": Group,
    "groupAvatar": File,
    "removed": User,
    "removedAvatar": File,
    "removedBy": User,
    "removedByAvatar": File,
    "seen (optional)": false
}
```