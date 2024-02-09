# Go Webhook Server

[![wakatime](https://wakatime.com/badge/user/bfa5a500-7b93-4deb-a695-4567ab9e77a8/project/018d7932-9a15-4a60-bb32-cbc13a35a2a9.svg)](https://wakatime.com/badge/user/bfa5a500-7b93-4deb-a695-4567ab9e77a8/project/018d7932-9a15-4a60-bb32-cbc13a35a2a9)

Build a webhook server that...

- receives webhooks
- handles duplicated webhooks
- idempotentally handles webhook retries - TODO
- offloads webhook processing to background jobs - TODO
- allow for failure and resumability for background jobs - TODO

## Goals

- Build a portfolio piece.
- Spend at least 100 hours on this project.
- Improve at one of my weaker languages (either golang or elixir). Felt that golang had a better market for me as an early career.
- Be technically challenging.

## Technology

- golang
- sqlite
- go-migrate
- hand rolled SQL (no ORMs for this project)
- [spf13/cobra](https://github.com/spf13/cobra)

## Acknowlegdements

Shouting out some resources that I've used whilst building this.

**GopherCon 2019: Mat Ryer - How I Write HTTP Web Services after Eight Years**
https://www.youtube.com/watch?v=rWBSMsLG8po

**GopherCon 2017: Mitchell Hashimoto - Advanced Testing with Go**
https://www.youtube.com/watch?v=8hQG7QlcLBk

## Project Structure

The application can be broken down into multiple sub-applications:

- `server`: The main application, the Webhook server. Accessed via the `server` CLI command.
- `sender-cli`: CLI app to send test webhook events to the server. Accessed via the `send` CLI command.
- `migrations-cli`: Applies schema changes to the database. It can be accessed via the `migrate-up` and `migrate-down` CLI commands.

### Packages

- `database`: Database utils such as opening a db connection.
- `logger`: Structured logging via uber-go/zap.
- `migration`: Applies database migrations. Used by the migrate-up & migrate-down CLI commands.
- `models`: Defines the struct types for the data models used in the database.
- `repository`: Implements CRUD operations for persisting data to the database.
- `sender`: Send a webhook event to the Webhook server. Used by the `send` CLI command.
- `server`: Webhook server. This is the main application in this project.
