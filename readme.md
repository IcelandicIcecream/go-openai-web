# Go OpenAI Web Server Implementation

### This is a super simple stateless implementation of a web server using the unofficial [Go client](https://github.com/sashabaranov/go-openai) for the OpenAI API. Uses SSE for streaming responses and the Go standard http library. This was a weekend project, so it's super bare-bones. Feel free to use it as a base for your own projects.

### Requirements

- [Go](https://go.dev/) version 1.22.0
- Make (optional)

### Usage

- Add your OpenAI API Key to .env.local and rename it to .env
- Run the `make start-server` OR `go run main.go` command

#### Recommended Front-end implementation

- I'm currently using [SSE.js](https://github.com/mpetazzoni/sse.js) for the Front-end SSE implementation. Hence, my `/openai` route accepts a post request with a JSON payload but returns
  an SSE streaming response.
