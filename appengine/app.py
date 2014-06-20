import webapp2
import hashlib

from google.appengine.ext import ndb

class Secret(ndb.Model):
  secret = ndb.StringProperty(required=True)


class DataModel(ndb.Model):
  data = ndb.StringProperty()
  created = ndb.DateTimeProperty(auto_now_add=True)


class PostJson(webapp2.RequestHandler):
  def post(self):
    # Check HMAC signature.
    signature = self.requests.get('signature')
    data = self.request.get('data')
    m = hashlib.sha256()
    m.update(data)
    secret = Secret.all().fetch(1)[0].secret
    m.update(data)
    m.update(secret)
    if m.hexdigest() != signature:
      self.error(403)
      return

    point = DataModel(data=data)
    point.put()
    self.response.write('OK')
    return
