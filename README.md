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
| **M**ove    | id, new week, new position  | Move a talk to a different slot |
| **N**ew     | id, name, type, description | Create a new talk               |
| **H**ide    | id                          | Hide a talk from the meeting    |
| **D**elete  | id                          | Delete a hidden talk            |
| **A**uth    | password                    | Authenticate client             |
