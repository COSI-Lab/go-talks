## Endpoints

| Request | Endpoint           | Desc                                             |
| :------ | :----------------- | :----------------------------------------------- |
| GET     | /                  | The talks homepage                               |
| GET     | /all               | Returns all former talks                         |
| GET     | /talks             | JSON of the currently visible talks              |
| GET     | /health            | Indicates how many active connections there are  |
| GET     | /ws                | Websocket endpoint                               |

## Websocket messages

| Type    | Arguments                   | Description                     |
| :-----: | :------------------------:  | :-----------------------------: |
| Move    | id, new week, new position  | Move a talk to a different slot |
| New     | id, name, type, description | Create a new talk               |
| Hide    | id                          | Hide a talk from the meeting    |
| Delete  | id                          | Delete a hidden talk            |
| Auth    | password                    | Authenticate client             |
