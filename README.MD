# psst-backend

HTTP server written in `go` for handling traffic from [psst-app](https://github.com/magicznykacpur/psst-app)

## requirements

- `go` version at least `1.24.1`
- `postgresql` you'll need one DB for running tests as well as one for the server itself
- `goose` for running migrations
- a copy of `.env.example` named `.env` with your environment variables set up correctly