# go-talks

Another rewrite of talks, this time in go.

go-talks (Or more commonly known as just Talks) is an app to manage talks at COSI meetings. It allows people to submit talks that they are planning on giving at upcoming meetings.

![preview](talkspreview.png)

## Features

- Live updating over websocket
- The page will load without Javascript
- Markdown support
- /posts directory supports rendering pre-written markdown files
- Image caching/proxying
- Create future talks and view historic talks
- Backed by SQLITE database

## Endpoints

| Request | Endpoint           | Desc                                             |
| :------ | :----------------- | :----------------------------------------------- |
| GET     | /                  | The talks homepage                               |
| GET     | /talks             | JSON of the talks for the next meeting           |
| GET     | /{week}            | The talks for the given week                     |
| GET     | /{week}/talks      | JSON of the talks for a week                     |
| GET     | /img/{id}          | Image proxy                                      |
| GET     | /health            | Indicates how many active connections there are  |
| GET     | /ws                | Websocket endpoint                               |
