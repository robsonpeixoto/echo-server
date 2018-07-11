import os
import pprint

from flask import Flask, jsonify, request

app = Flask(__name__)
app.config['JSONIFY_PRETTYPRINT_REGULAR'] = True

app_name = os.getenv('APP_NAME')

ALL_METHODS = ['GET', 'HEAD', 'POST', 'PUT', 'DELETE', 'OPTIONS']


@app.route('/', defaults={'path': ''}, methods=ALL_METHODS)
@app.route('/<path:path>', methods=ALL_METHODS)
def index(path):
    data = {
        'path': request.path,
        'method': request.method,
        'headers': list(request.headers.items()),
        'form': list(request.form.items()),
        'args': list(request.args.items()),
        'remote': {
            'address': request.environ['REMOTE_ADDR'],
            'port': request.environ['REMOTE_PORT'],
        },
        'content_type': request.content_type,
        'files': [(f[0], f[1].filename) for f in request.files.items()],
        'json': request.json,
        'raw-data': str(request.data)
    }

    if app_name:
        data['APP_NAME'] = app_name
    app.logger.info('\n' + pprint.pformat(data))
    return jsonify(data)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
