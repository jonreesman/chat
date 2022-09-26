# chat-app
This chat application was a learning tool to experiment with websockets and concurrency. I started it as a long-term project to make my own implementation of the matrix.org protocols based on [element.io](https://element.io/). In its present state, it is a simple chat room application with user authentication that secures critical endpoints with same-site, HTTP Only JWTs to prevent XSS.

The next step in the project will be implementing end-to-end encryption. My current plan is to utilize Go compiled to WASM in the frontend for all encryption and decryption.

It has an accompanying Next.js based frontend located [here.](https://github.com/jonreesman/chat-next) A live demo does not yet exist.

## Features
- user authentication
- custom user avatars
- user nicknames
- persisted chat room messages
- chat room creation/deletion

## Set-Up
1. Use `docker-compose.yml` to "compose up" the Postgres container.
    - ENV variables: See the included [.env.sample](https://github.com/jonreesman/chat/blob/master/.env.sample).
2. Run `go build`
3. Run `./chat`
4. That's all it takes to start up the backend!

## Notes:
- The app utilizes GORM, so there is minimal configuration needed to get PostgreSQL running properly!
- By default, there is no user accounts made upon creation. To create one, utilize the `createClient` endpoint referenced [here](https://github.com/jonreesman/chat/blob/master/router/router.go#L53) in the `router/router.go` package, and defined [here](https://github.com/jonreesman/chat/blob/master/handler/client.go#L71) in the `handler/client.go` package.
- All user passwords are stored in the DB as a hash.
- The `.env.sample` and included `docker-compose.yml` will work together out of the box. You will want to set your own `SECRET` variable though!


