## Endpoints

| Request | Endpoint           | Desc                                             |
| :------ | :----------------- | :----------------------------------------------- |
| GET     | /                  | The talks homepage                               |
| GET     | /all               | Returns all former talks                         |
| GET     | /talks             | JSON of the currently visible talks              |
| GET     | /health            | Indicates how many active connections there are  |
| POST    | /register          | Registers a new client for live updates          |
| POST    | /authenticate      | authenticates a client                           |
| GET     | /ws/{id}           | Websocket endpoint                               |
