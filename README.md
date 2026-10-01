# echo-server

A simple HTTP server that echoes request details back as JSON. Useful to test and debug HTTP clients, proxies and ingress rules.

## Docker Registries

Images are published for `linux/amd64` and `linux/arm64`.

- [Docker Hub](https://hub.docker.com/r/robsonpeixoto/echo-server)
- [AWS ECR Public](https://gallery.ecr.aws/v2m3p9l8/robsonpeixoto/echo-server)
- [Github Docker Registry](https://github.com/robsonpeixoto/echo-server/pkgs/container/echo-server)

## Usage

```sh
docker run --rm -p 5000:5000 -e APP_NAME=robinho robsonpeixoto/echo-server
```

Or without Docker:

```sh
go run .
```

In another terminal:

```sh
❯ curl -s 'localhost:5000/hello?name=robinho' | jq
{
  "host": "localhost:5000",
  "proto": "HTTP/1.1",
  "content_length": 0,
  "headers": {
    "User-Agent": ["curl/8.18.0"],
    "Accept": ["*/*"]
  },
  "form": {},
  "query": {
    "name": ["robinho"]
  },
  "remote": {
    "address": "172.17.0.1",
    "port": "34508"
  },
  "path": "/hello",
  "method": "GET",
  "extras": {
    "app_name": "robinho"
  }
}
```

JSON bodies are echoed under `json`, and URL-encoded forms under `form`:

```sh
❯ curl -s -H 'Content-Type: application/json' -d '{"message":"hi"}' localhost:5000 | jq .json
{
  "message": "hi"
}

❯ curl -s -d 'user=robinho' localhost:5000 | jq .form
{
  "user": ["robinho"]
}
```

A request with `Content-Type: application/json` and an invalid body returns `400`.

## Configuration

| Env / flag  | Default | Description                                   |
|-------------|---------|-----------------------------------------------|
| `PORT`      | `5000`  | Port to listen on                             |
| `APP_NAME`  |         | Returned as `extras.app_name`                 |
| `SHOW_ENVS` |         | When `1`, returns all env vars in `extras.envs` |
| `-version`  |         | Print version and commit, then exit           |

> [!WARNING]
> `SHOW_ENVS=1` exposes every environment variable, including secrets, to anyone who can reach the server.
