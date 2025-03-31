# psst-backend

HTTP server written in `go` for handling traffic for a communicator application.

## Why
I've always wondered how apps like `discord` and `slack` work, so i've set out to build
something similiar myself. Although the end result might be not nearly close to the major
competitors, I've learned a lot along the way. How to handle users and their authentication/authorization,
how to expose a server on a local network, how to handle notiyfing and receiving messages via websockets.

## Quick start

### requirements

- `go` version at least `1.24.1`
- `postgresql` you'll need one DB for running tests as well as one for the server itself
- `goose` for running migrations
- a copy of `.env.example` named `.env` with your environment variables set up correctly

### install and run
- run `go install` to download necessary packages
- create your local database using `psql` or `pg-admin`, you can name it however you like just remember to set up the correct url in the `.env` file
- run migrations with `goose`, there's a `shell` script I've wrote in `sql/schema`
- run the app with `go run .` from the root directory of your local copy, you should see `starting psst server at` log

## Usage
`psst-backend` exposes a few useful endpoints on a provided port of your local machine

- `GET /api/users` returns all the users
- `GET /api/users/{id}` returns the user with id
- `POST /api/users` creates new user

- `GET /api/users/me` returns the user with a valid JWT token
- `GET /api/users/to-chat-with` returns users a chat can be created with

- `POST /api/login` returns a signed JWT token provided valid email and password

- `POST /api/messages` creates a new message in a chat
- `GET /api/messages/{chat_id}` returns all the messages from a chat
- `DELETE /api/messages` deletes a message from a chat

- `POST /api/chats` creates a new chat
- `GET /api/chats` returns all logged in users chats
- `DELETE /api/chats/{id}` deletes a chat

- `/ws` creates and registers a websocket client


### about

`psst-backend` exposes several RESTful endpoints that can be used to communicate
with a database. It also exposes a `websocket` endpoint that you can connect to 
and receive messages about new chats and messages.