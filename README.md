# echo-server

```sh
❯ docker run --rm -p 5000:5000 -e PORT=5000 -e APP_NAME=robinho -e SHOW_ENVS=1 robsonpeixoto/echo-server

STARTING
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

Will show the log:

```
{"time":"2023-08-09T10:42:10.000541568Z","level":"INFO","msg":"","response":{"headers":{"Accept":["*/*"],"Accept-Encoding":["gzip, deflate"],"Connection":["keep-alive"],"User-Agent":["HTTPie/3.2.2"]},"form":null,"query":{},"remote":{"address":"192.168.215.1","port":"51628"},"path":"/","method":"GET","extras":{"envs":{"APP_NAME":"robinho","HOME":"/root","HOSTNAME":"b8abc5c78e4f","PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","PORT":"5000","SHOW_ENVS":"1"},"app_name":"robinho"}}}
```
