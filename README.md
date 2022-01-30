## Development

**Auto reload:**
```
gin --all -d src -b main -a 3001 -i run
```

## Endpoints

| Request | Endpoint           | Desc                                             |
| :------ | :----------------- | :----------------------------------------------- |
| GET     | /                  | The talks homepage                               |
| GET     | /all               | Returns all former talks                         |
| GET     | /talks             | JSON of the currently visible talks              |
| GET     | /health            | Indicates how many active connections there are  |
| GET     | /ws                | Websocket endpoint                               |

**Create a new talk**
```json
{
    "type": 0,
    "name": "mirror team",
    "talktype": 0, // forum topic
    "description": "new drives for mirror",
    "week": 20220126, // YYYYMMDD optional, what day to add the talk to. empty for current week otherwise must be in the future and a Wednesday
}
```

**Hide a talk**
```json
{
    "type": 1,
    "id": 10, // talk id
}
```

**Delete a talk**
```json
{
    "type": 2,
    "id": 10, // talk id
}
```

TODO: MOVE TALKS

**Authenticate** (must be sent before you can change talk states)
```json
{
    "type": 4,
    "password": "password",
}
```
If the password is correct the server sends back 
```json
{
    "type": 4,
    "status": 200,
}
```
otherwise
```json
{
    "type": 4,
    "status": 403,
}
```