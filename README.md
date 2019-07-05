# echo-server

```sh
❯ docker run --rm -p 5000:5000 -e PORT=5000 -e APP_NAME=robinho -e SHOW_ENVS=1 robsonpeixoto/echo-server

STARTING
```

In other terminal execute the command:

```sh
❯ http localhost:5000
HTTP/1.1 200 OK
Content-Length: 992
Content-Type: application/json
Date: Fri, 05 Jul 2019 00:25:32 GMT
Server: waitress

{
    "args": [],
    "content-type": null,
    "extras": {
        "app_name": "robinho",
        "envs": {
            "APP_NAME": "robinho",
            "GPG_KEY": "0D96DF4D4110E5C43FBFB17F2D347EA6AA65421D",
            "HOME": "/root",
            "HOSTNAME": "0a9e466aaccc",
            "LANG": "C.UTF-8",
            "PATH": "/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "PORT": "5000",
            "PWD": "/usr/src/app",
            "PYTHON_PIP_VERSION": "19.1.1",
            "PYTHON_VERSION": "3.7.3",
            "SHOW_ENVS": "1"
        }
    },
    "files": [],
    "form": [],
    "headers": [
        [
            "Host",
            "localhost:5000"
        ],
        [
            "User-Agent",
            "HTTPie/1.0.2"
        ],
        [
            "Accept-Encoding",
            "gzip, deflate"
        ],
        [
            "Accept",
            "*/*"
        ],
        [
            "Connection",
            "keep-alive"
        ]
    ],
    "json": null,
    "method": "GET",
    "path": "/",
    "raw-data": "b''",
    "remote": {
        "address": "172.17.0.1",
        "port": "36610"
    }
}
```

Will show the log:

```
[2019-07-05 00:25:32,018] INFO in run:
{'args': [],
 'content-type': None,
 'extras': {'app_name': 'robinho',
            'envs': {'APP_NAME': 'robinho',
                     'GPG_KEY': '0D96DF4D4110E5C43FBFB17F2D347EA6AA65421D',
                     'HOME': '/root',
                     'HOSTNAME': '0a9e466aaccc',
                     'LANG': 'C.UTF-8',
                     'PATH': '/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
                     'PORT': '5000',
                     'PWD': '/usr/src/app',
                     'PYTHON_PIP_VERSION': '19.1.1',
                     'PYTHON_VERSION': '3.7.3',
                     'SHOW_ENVS': '1'}},
 'files': [],
 'form': [],
 'headers': [('Host', 'localhost:5000'),
             ('User-Agent', 'HTTPie/1.0.2'),
             ('Accept-Encoding', 'gzip, deflate'),
             ('Accept', '*/*'),
             ('Connection', 'keep-alive')],
 'json': None,
 'method': 'GET',
 'path': '/',
 'raw-data': "b''",
 'remote': {'address': '172.17.0.1', 'port': '36610'}}
```
