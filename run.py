from flask import Flask, request
from gevent.wsgi import WSGIServer
import pprint

app = Flask(__name__)

ALL_METHODS = ['GET', 'HEAD', 'POST', 'PUT', 'DELETE', 'OPTIONS']

@app.route('/', defaults={'path': ''}, methods=ALL_METHODS)
@app.route('/<path:path>', methods=ALL_METHODS)
def index(path):
    data = dict(
        path = request.path,
        method = request.method,
        headers = list(request.headers.items()),
        form = list(request.form.items()),
        args = list(request.args.items()),
        files = list(request.files.items()),
        json = request.json,
        content_type = request.content_type,
    )
    app.logger.error('\n' + pprint.pformat(data))
    return "Ok!"

if __name__ == '__main__':
    http_server = WSGIServer(('', 5000), app)
    http_server.serve_forever()