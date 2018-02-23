import pprint

from flask import Flask, jsonify, request
from gevent.wsgi import WSGIServer

app = Flask(__name__)
app.config['DEBUG'] = True

ALL_METHODS = ['GET', 'HEAD', 'POST', 'PUT', 'DELETE', 'OPTIONS']


@app.route('/', defaults={'path': ''}, methods=ALL_METHODS)
@app.route('/<path:path>', methods=ALL_METHODS)
def index(path):
    data = dict(
        path=request.path,
        method=request.method,
        headers=list(request.headers.items()),
        form=list(request.form.items()),
        args=list(request.args.items()),
        remote=dict(
            address=request.environ['REMOTE_ADDR'],
            port=request.environ['REMOTE_PORT'],
        ),
        content_type=request.content_type,
        files=[(f[0], f[1].filename) for f in request.files.items()],
        json=request.json,
    )
    app.logger.info('\n' + pprint.pformat(data))
    return jsonify(data)


if __name__ == '__main__':
    port = 5000
    http_server = WSGIServer(('', port), app)
    app.logger.info('RUNNING on port', port)
    http_server.serve_forever()
