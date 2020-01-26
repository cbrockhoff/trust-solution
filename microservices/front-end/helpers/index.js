(function (){
  'use strict';

  var endpoints = require("../api/endpoints");
  var request = require("request");
  var helpers = {};

  helpers.errorHandler = function(err, req, res, next) {
    var ret = {
      message: err.message,
      error:   err
    };
    res.
      status(err.status || 500).
      send(ret);
  };

  helpers.sessionMiddleware = function(err, req, res, next) {
    if(!req.cookies.logged_in) {
      res.session.customerId = null;
    }
  };

  /* Responds with the given body and status 200 OK  */
  helpers.respondSuccessBody = function(res, body) {
    helpers.respondStatusBody(res, 200, body);
  }

  helpers.respondStatusBody = function(res, statusCode, body) {
    res.writeHeader(statusCode);
    res.write(body);
    res.end();
  }

  helpers.respondStatus = function(res, statusCode) {
    res.writeHeader(statusCode);
    res.end();
  }

  helpers.rewriteSlash = function(req, res, next) {
   if(req.url.substr(-1) == '/' && req.url.length > 1)
       res.redirect(301, req.url.slice(0, -1));
   else
       next();
  }

  helpers.getServiceToken = function(callback) {
    const username = 'front-end'
    const password = 'front-end_password'

    try {
      if (typeof helpers.fetchTokenAt === 'undefined' || Date.now() >= helpers.fetchTokenAt) {
        var options = {
          headers: {
            'Authorization': "Basic " + new Buffer(username + ":" + password).toString('base64')
          },
          uri: endpoints.tokensUrl,
          method: 'GET'
        };
        
        request(options, function(err, resp, body) {
          if (err) {
            console.log("Error fetching service token: " + err);
          } else {
            console.log("Received response from tokenservice with body " + JSON.stringify(body))
            helpers.serviceToken = JSON.parse(body).Token;
            helpers.fetchTokenAt = Date.now() + 4.5 * 60 * 1000; // 4.5 minutes * 60 seconds * 1000 milliseconds... Should really get this from the token. Oh well
            callback(helpers.serviceToken);
          }
        });
      } else {
        console.log("Returning saved Service-Token, fetchTokenAt: " + helpers.fetchTokenAt);
        callback(helpers.serviceToken);
      }
    } catch (err) {
      console.log("Exception while fetching service token: " + err);
      callback("");
    }
  }

  helpers.simpleHttpRequest = function(url, userToken, res, next) {
    helpers.getServiceToken(serviceToken => {
      var options = {
        headers: {
          'User-Token' : userToken ? userToken : "",
          'Service-Token' : serviceToken
        },
        uri: url,
        method: 'GET',
      };
      
      if (!userToken || userToken === "") {
        console.log("User-Token is empty");
      }

      console.log("Making simple request: " + options.method + " " + options.uri);

      try {
        request(options, function(err, resp, body) {
          if (err) {
            console.log("Error: " + err);
            return next(err);
          }
          helpers.respondSuccessBody(res, body);
        });
      } catch (err) {
        console.log("Error making request: " + err);
      }
    })
  }

  helpers.getCustomerId = function(req, env) {
    var logged_in = req.cookies.logged_in;

    if (env == "development" && req.query.custId != null) {
      return req.query.custId;
    }

    if (!logged_in) {
      if (!req.session.id) {
        throw new Error("User not logged in.");
      }
      console.log("User not logged in, using session ID " + req.session.id);
      
      return req.session.id;
    }

    console.log("User logged in, with customer ID " + logged_in);
    return logged_in;
  }
  module.exports = helpers;
}());
