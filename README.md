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

