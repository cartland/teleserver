import webapp2
import hashlib
import json

from google.appengine.ext import ndb

class Secret(ndb.Model):
  secret = ndb.StringProperty(required=True)


class DataModel(ndb.Model):
  data = ndb.StringProperty()
  created = ndb.DateTimeProperty(auto_now_add=True)


class SecretPage(webapp2.RequestHandler):
  def post(self):
    new_secret = self.request.get('new_secret')
    secret_o = Secret.query().get()
    if not secret_o:
      secret_o = Secret(secret='default')
      secret_o.put()
    if new_secret:
      secret_o.secret = new_secret
      secret_o.put()

    self.response.write('Secret: %s' % secret_o.secret)

  def get(self):
    return self.post()


class PostJson(webapp2.RequestHandler):
  def post(self):
    # Check HMAC signature.
    signature = self.request.get('signature')
    data = self.request.get('data')
    m = hashlib.sha256()
    secret_model = Secret.query().get()
    if not secret_model:
      self.response.write('No secret set')
      self.error(501)
      return
    secret = secret_model.secret
    m.update(data)
    m.update(secret)
    if m.hexdigest() != signature:
      self.error(403)
      return

    point = DataModel(data=data)
    point.put()
    self.response.write('OK')
    return

  def get(self):
    return self.post()


class FetchData(webapp2.RequestHandler):
  def get(self):
    limit = self.request.get('limit', 10)
    result = []
    for data in DataModel.query().order(-DataModel.created).fetch(limit):
      result.append(data.data)
    self.response.write(json.dumps(result))


app = webapp2.WSGIApplication([
    ('/admin/secret', SecretPage),
    ('/fetch_data', FetchData),
    ('/post', PostJson),
], debug=True)
