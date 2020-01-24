(function (){
  'use strict';

  var express   = require("express")
    , request   = require("request")
    , endpoints = require("../endpoints")
    , helpers   = require("../../helpers")
    , app       = express()

  app.get("/catalogue/images*", function (req, res, next) {
    helpers.getServiceToken(serviceToken => {
      var options = {
        headers: {
          "User-Token": req.cookies.user_token,
          "Service-Token": serviceToken
        },
        uri: endpoints.catalogueUrl + req.url.toString()
      }
      request(options)
          .on('error', function(e) { next(e); })
          .pipe(res);
    })
  });

  app.get("/catalogue*", function (req, res, next) {
    helpers.simpleHttpRequest(endpoints.catalogueUrl + req.url.toString(), req.cookies.user_token, res, next);
  });

  app.get("/tags", function(req, res, next) {
    helpers.simpleHttpRequest(endpoints.tagsUrl, req.cookies.user_token, res, next);
  });

  module.exports = app;
}());
