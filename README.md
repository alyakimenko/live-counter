I decided to use SSE (Server-Sent Events), instead of WebSockets.

There are several reasons for this:
1. I think that SSE is a better way to solve the task.
2. I've never tried it before, and this was my chance to try.

# Installing

`git clone https://github.com/alyakimenko/live-counter`

# Running

To run the server, do

`make run`

or if you don't have the `make`, you can just

`go run ./cmd/...`

Then point your browser to `http://localhost:8080`