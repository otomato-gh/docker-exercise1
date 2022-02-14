from flask import Flask, request, render_template, redirect, url_for
import requests
from requests.adapters import HTTPAdapter
import json
from werkzeug.utils import secure_filename
from flask_wtf import FlaskForm
from wtforms import StringField, PasswordField, BooleanField, DecimalField, FileField, SubmitField
from wtforms.validators import DataRequired, InputRequired, Email, Length
import os, time


app = Flask(__name__)
app.config.from_object('version')
SECRET_KEY = os.urandom(32)
app.config['SECRET_KEY'] = SECRET_KEY


from flask_bootstrap import Bootstrap
Bootstrap(app)


class MyForm(FlaskForm):
    name = StringField('title', validators=[DataRequired()])
    type = StringField('Desc', validators=[DataRequired()])
    price = StringField('Content', validators=[DataRequired()])
    picture = FileField('picture')
    submit = SubmitField(label='Add')


@app.route('/healthz')
def healthz():
    return 'Healthy'

@app.route('/version')
def version():
    return app.config['VERSION']


@app.route('/login', methods=['POST'])
def login():
    user = request.values.get('username')
    response = app.make_response(redirect(request.referrer))
    response.set_cookie('user', user)
    return response


@app.route('/logout', methods=['GET'])
def logout():
    response = app.make_response(redirect(request.referrer))
    response.set_cookie('user', '', expires=0)
    return response

@app.route('/submit', methods=('GET', 'POST'))
def submit():
    print("i'm in submit")
    print("Title"+request.values.get('Title'))
    response = app.make_response(redirect(request.referrer))
  
    api_port = os.getenv("API_PORT", "8888")
    api_url = "http://api:" + api_port
    r = requests.post(api_url + "/post",
                    data = { "Title": request.values.get('Title'),
                            "Desc": request.values.get('Desc'),
                            "Content": request.values.get('Content')
                            })
    print(r.text)
    return response

@app.route('/')
def home():
    user = request.cookies.get("user", "")
    form = MyForm()

    api_port = os.getenv("API_PORT", "8888")
    api_url = "http://api:" + api_port
    r = requests.get(api_url +"/posts")
    if r.status_code != 500:
        print(r.text)
        posts = json.loads(r.text)
        print(posts)
    return render_template(
        'front.html',
        posts=posts,
        user=user,
        form=form
    )


def run_app():
    s = requests.Session()
    app.run(host='0.0.0.0', port=os.getenv("PORT", "8888"), debug=True)

if __name__ == '__main__':
    run_app()
