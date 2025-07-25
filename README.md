# echo-server

Just a simple http application to test and debug HTTP requests.

## Docker Registries

- [Docker Hub](https://hub.docker.com/r/robsonpeixoto/echo-server)
- [AWS ECR Public](https://gallery.ecr.aws/v2m3p9l8/robsonpeixoto/echo-server)
- [Github Docker Registry](https://github.com/robsonpeixoto/echo-server/pkgs/container/echo-server)

## Usage

```sh
❯ docker run --rm -p 5000:5000 -e PORT=5000 -e APP_NAME=robinho -e SHOW_ENVS=1 robsonpeixoto/echo-server
```

In other terminal execute the command:

```sh
❯ http localhost:5000

HTTP/1.1 200 OK
Content-Length: 428
Content-Type: application/json
Date: Wed, 09 Aug 2023 10:42:10 GMT

{
    "extras": {
        "app_name": "robinho",
        "envs": {
            "APP_NAME": "robinho",
            "HOME": "/root",
            "HOSTNAME": "b8abc5c78e4f",
            "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "PORT": "5000",
            "SHOW_ENVS": "1"
        }
    },
    "form": null,
    "headers": {
        "Accept": [
            "*/*"
        ],
        "Accept-Encoding": [
            "gzip, deflate"
        ],
        "Connection": [
            "keep-alive"
        ],
        "User-Agent": [
            "HTTPie/3.2.2"
        ]
    },
    "method": "GET",
    "path": "/",
    "query": {},
    "remote": {
        "address": "192.168.215.1",
        "port": "51628"
    }
}

```
