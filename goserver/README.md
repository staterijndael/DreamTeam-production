## Structures

### User

```json
{
    "id": 0,
    "nickname": "",
    "score": Score,
    "teel": 0,
    "experience": 0,
    "(optional) firstName": "",
    "(optional) lastName": "",
    "(optional) bdate": 0,
    "(optional) aboutMe": "",
    "(optional) email": "",
    "(optional) phone": ""
}
```

#### http

##### PUT /user

```json
{
    "phone": "",
    "password": ""
}
```

returns

```json
{
    "id": 0
}
```

##### POST /auth

```json
{
    "nickname | phone | email": 0,
    "password": ""
}
```

returns

```json
{
    "id": 0,
    "token": ""
}
```

### Organization

```json
{
    "id": 0,
    "title": "",
    "nickname": "",
    "description": "",
    "director": 0,
    "associatedUsers": [0],
    "date": 0
}
```

### File

```json
{
    "id": 0,
    "checksum": "",
    "size": 0,
    "content": ""
}
```

### FileMetaInfo

```json
{
    "id": 0,
    "checksum": "",
    "size": 0
}
```

### FileContent

```json
{
    "content": ""
}
```

### Score

```json
{
    "efficiency": 0,
    "loyalty": 0,
    "professionalism": 0,
    "discipline": 0,
    "rang": ""
}
```

### AdministratedOrganizations
```json
{
  "director": [Organization],
  "linked": [Organization]
}
```

### RatingPreferences

```json
{
    "id": 0,
    "start": 0,
    "end": 0,
    "orgID": 0
}
```

### GroupRatingEvent

```json
{
    "id": 0,
    "group": "",
    "date": 0,
    "preferences": RatingPreferences
}
```

### AverageGroupRatingOfUser

```json
{
    "event": 0,
    "user": 0,
    "score": Score
}
```

### AverageOrganizationRatingOfUser

```json
{
    "id": 0,
    "user": 0,
    "date": 0,
    "score": Score
}
```

### Container

```json
{
    "type": "",
    "data": {}
}
```

### Community

```json
{
    "id": 0,
    "members": [0]
}
```

### GroupCreationRequest

```json
{
  "id": 0,
  "status": "",
  "initiator": 0,
  "acceptor (optional)": 0,
  "groupInfo": {
    "parent (optional)": "",
    "title": "",
    "description": ""
  },
  "organization": 0
}
```

### Group

```json
{
    "id": "",
    "title": "",
    "description": "",
    "children": ["groupsIDs"],
    "parent (optional)": "",
    "org": 0,
    "community": 0,
    "admin": 0,
    "creator": 0,
    "status": "",
    "date": 0
}
```

### GroupInvitation

```json
{
    "id": 0,
    "status": "",
    "receiver": 0,
    "initiator": 0,
    "group": "",
    "date": 0,
    "accepter (optional)": 0,
    "repliedAt (optional)": 0
}
```

### State

```json
{
    "groupAdministrating": {
        "usersAddedToGroup": [{
            "user": 0,
            "group": ""
        }],
        "usersLeft": [{
            "user": 0,
            "group": ""
        }],
        "byOrganizator": {
            "usersRemoved": [{
                "user": 0,
                "group": "",
                "organizator": 0
            }],
            "groupInvited": [{
                "invitation": 0
            }],
            "creationAccepted":{
                "group": "",
                "by": 0
            }
        }
    },
    "groups": {
        "kickedOut": [""],
        "invited": [{
            "initiator": 0,
            "group":""
        }],
        "added": {
            "invitation": 0
        },
        "creationAccepted":{
            "group": "",
            "by": 0
        },
        "groupDeleted": {}
    },
    "organizating": {
        "usersAddedToGroup": [{
            "user": 0,
            "group": ""
        }],
        "created": {
            "group": "",
            "creator": 0
        },
        "groupInvited": [{
            "initiator": 0,
            "group":"",
            "receiver": 0
        }],
        "newUserInGroup": {
            "invitation": 0
        },
        "usersRemovedByOrganizator": [{
            "user": 0,
            "group": "",
            "organizator": 0
        }],
        "usersRemovedByAdmin": [{
            "user": 0,
            "group": "",
            "admin": 0
        }],
        "usersLeft": [{
            "user": 0,
            "group": ""
        }]
    }
}
```

### Rating

#### RatingQuestions

```json
{
	"discipline": [RatingQuestion],
	"professionalism": [RatingQuestion],
	"efficiency": [RatingQuestion],
	"loyalty": [RatingQuestion]
}
```

#### RatingQuestion

```json
{
	"text": "",
	"required": true,
	"options (optional)": [RatingQuestionOption]
}
```

#### RatingQuestionOption

```json
{
	"text": "",
	"cost": 0
}
```