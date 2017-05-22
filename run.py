from flask import Flask, request, jsonify
from gevent.wsgi import WSGIServer
import pprint

app = Flask(__name__)
app.config['DEBUG'] = True

ALL_METHODS = ['GET', 'HEAD', 'POST', 'PUT', 'DELETE', 'OPTIONS']

def get_file_args(files):
    return [(f[0], f[1].filename) for f in files]

@app.route('/', defaults={'path': ''}, methods=ALL_METHODS)
@app.route('/<path:path>', methods=ALL_METHODS)
def index(path):
    data = dict(
        path = request.path,
        method = request.method,
        headers = list(request.headers.items()),
        form = list(request.form.items()),
        args = list(request.args.items()),
        files = list(get_file_args(request.files.items())),
        json = request.json,
        content_type = request.content_type,
    )
    app.logger.info('\n' + pprint.pformat(data))
    return jsonify(data)

if __name__ == '__main__':
    http_server = WSGIServer(('', 5000), app)
    app.logger.info('RUNNING')
    http_server.serve_forever()

