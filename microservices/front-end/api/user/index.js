(function() {
    'use strict';

    var async = require("async"), express = require("express"), request = require("request"), endpoints = require("../endpoints"), helpers = require("../../helpers"), app = express(), cookie_name = "logged_in", cookie_token_name = "user_token"

    app.get("/customers/:id", function(req, res, next) {
        helpers.simpleHttpRequest(endpoints.customersUrl + "/" + req.cookies.logged_in, req.cookies.user_token, res, next);
    });
    app.get("/cards/:id", function(req, res, next) {
        helpers.simpleHttpRequest(endpoints.cardsUrl + "/" + req.params.id, req.cookies.user_token, res, next);
    });
    app.get("/customers", function(req, res, next) {
        helpers.simpleHttpRequest(endpoints.customersUrl, req.cookies.user_token, res, next);
    });
    app.get("/addresses", function(req, res, next) {
        helpers.simpleHttpRequest(endpoints.addressUrl, req.cookies.user_token, res, next);
    });
    app.get("/cards", function(req, res, next) {
        helpers.simpleHttpRequest(endpoints.cardsUrl, req.cookies.user_token, res, next);
    });

    // Create Customer - TO BE USED FOR TESTING ONLY (for now)
    app.post("/customers", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.customersUrl,
                method: 'POST',
                json: true,
                body: req.body
            };

            console.log("Posting Customer: " + JSON.stringify(req.body));

            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    app.post("/addresses", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            req.body.userID = helpers.getCustomerId(req, app.get("env"));

            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.addressUrl,
                method: 'POST',
                json: true,
                body: req.body
            };
            console.log("Posting Address: " + JSON.stringify(req.body));
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    app.get("/card", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            var custId = helpers.getCustomerId(req, app.get("env"));
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.customersUrl + '/' + custId + '/cards',
                method: 'GET',
            };
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                var data = JSON.parse(body);
                if (data.status_code !== 500 && data._embedded.card.length !== 0 ) {
                    var resp = {
                        "number": data._embedded.card[0].longNum.slice(-4)
                    };
                    return helpers.respondSuccessBody(res, JSON.stringify(resp));
                }
                return helpers.respondSuccessBody(res, JSON.stringify({"status_code": 500}));
            }.bind({
                res: res
            }));
        })
    });

    app.get("/address", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            var custId = helpers.getCustomerId(req, app.get("env"));
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.customersUrl + '/' + custId + '/addresses',
                method: 'GET',
            };
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                var data = JSON.parse(body);
                if (data.status_code !== 500 && data._embedded.address.length !== 0 ) {
                    var resp = data._embedded.address[0];
                    return helpers.respondSuccessBody(res, JSON.stringify(resp));
                }
                return helpers.respondSuccessBody(res, JSON.stringify({"status_code": 500}));
            }.bind({
                res: res
            }));
        })
    });

    app.post("/cards", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            req.body.userID = helpers.getCustomerId(req, app.get("env"));

            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.cardsUrl,
                method: 'POST',
                json: true,
                body: req.body
            };
            console.log("Posting Card: " + JSON.stringify(req.body));
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    // Delete Customer - TO BE USED FOR TESTING ONLY (for now)
    app.delete("/customers/:id", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            console.log("Deleting Customer " + req.params.id);
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.customersUrl + "/" + req.params.id,
                method: 'DELETE'
            };
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    // Delete Address - TO BE USED FOR TESTING ONLY (for now)
    app.delete("/addresses/:id", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            console.log("Deleting Address " + req.params.id);
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.addressUrl + "/" + req.params.id,
                method: 'DELETE'
            };
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    // Delete Card - TO BE USED FOR TESTING ONLY (for now)
    app.delete("/cards/:id", function(req, res, next) {
        helpers.getServiceToken(serviceToken => {
            console.log("Deleting Card " + req.params.id);
            var options = {
                headers: {
                    'User-Token': req.cookies.user_token,
                    'Service-Token': serviceToken
                },
                uri: endpoints.cardsUrl + "/" + req.params.id,
                method: 'DELETE'
            };
            request(options, function(error, response, body) {
                if (error) {
                    return next(error);
                }
                helpers.respondSuccessBody(res, JSON.stringify(body));
            }.bind({
                res: res
            }));
        })
    });

    app.post("/register", function(req, res, next) {

        console.log("Posting Customer: " + JSON.stringify(req.body));

        async.waterfall([
                function(callback) {
                    helpers.getServiceToken(serviceToken => {
                        var options = {
                            headers: {
                                'Service-Token': serviceToken
                            },
                            uri: endpoints.registerUrl,
                            method: 'POST',
                            json: true,
                            body: req.body
                        };
                        request(options, function(error, response, body) {
                            if (error !== null ) {
                                callback(error);
                                return;
                            }
                            if (response.statusCode == 200 && body != null && body != "") {
                                if (body.error) {
                                    callback(body.error);
                                    return;
                                }
                                console.log(body);
                                var customerId = body.id;
                                var token = body.token;
                                req.session.customerId = customerId;
                                req.session.token = token;
                                callback(null, customerId, token);
                                return;
                            }
                            console.log(response.statusCode);
                            callback(true);
                        });
                    })
                },
                function(custId, token, callback) {
                    helpers.getServiceToken(serviceToken => {
                        var sessionId = req.session.id;
                        console.log("Merging carts for customer id: " + custId + " and session id: " + sessionId);

                        var options = {
                            headers: {
                                'User-Token' : token,
                                'Service-Token': serviceToken
                            },
                            uri: endpoints.cartsUrl + "/" + custId + "/merge" + "?sessionId=" + sessionId,
                            method: 'GET'
                        };
                        request(options, function(error, response, body) {
                            if (error) {
                                if(callback) callback(error);
                                return;
                            }
                            console.log('Carts merged.');
                            if(callback) callback(null, custId, token);
                        });
                    })
                }
            ],
            function(err, custId, token) {
                if (err) {
                    console.log("Error with log in: " + err);
                    res.status(500);
                    res.end();
                    return;
                }
                console.log("set cookie" + custId);
                res.status(200);
                res.cookie(cookie_token_name, token, {maxAge: 3600000});
                res.cookie(cookie_name, custId, {
                    maxAge: 3600000
                }).send({id: custId});
                console.log("Sent cookies.");
                res.end();
                return;
            }
        );
    });

    app.get("/login", function(req, res, next) {
        console.log("Received login request");

        async.waterfall([
                function(callback) {
                    helpers.getServiceToken(serviceToken => {
                        var options = {
                            headers: {
                                'Authorization': req.get('Authorization'),
                                'Service-Token': serviceToken
                            },
                            uri: endpoints.loginUrl
                        };
                        request(options, function(error, response, body) {
                            if (error) {
                                callback(error);
                                return;
                            }
                            if (response.statusCode == 200 && body != null && body != "") {
                                var parsedBody = JSON.parse(body);
                                var customerId = parsedBody.user.id;
                                console.log("User logged in as: " + customerId);
                                var token = parsedBody.token;
                                req.session.customerId = customerId;
                                req.session.token = token;
                                callback(null, customerId, token);
                                return;
                            }
                            console.log(response.statusCode);
                            callback(true);
                        });
                    })
                },
                function(custId, token, callback) {
                    helpers.getServiceToken(serviceToken => {
                        var sessionId = req.session.id;
                        console.log("Merging carts for customer id: " + custId + " and session id: " + sessionId);

                        var options = {
                            headers: {
                                'User-Token': token,
                                'Service-Token': serviceToken
                            },
                            uri: endpoints.cartsUrl + "/" + custId + "/merge" + "?sessionId=" + sessionId,
                            method: 'GET',
                        };
                        request(options, function(error, response, body) {
                            if (error) {
                                // if cart fails just log it, it prevenst login
                                console.log(error);
                                //return;
                            }
                            console.log('Carts merged.');
                            callback(null, custId);
                        });
                    })
                }
            ],
            function(err, custId) {
                if (err) {
                    console.log("Error with log in: " + err);
                    res.status(401);
                    res.end();
                    return;
                }
                res.status(200);
                res.cookie(cookie_token_name, req.session.token, {maxAge: 3600000});
                res.cookie(cookie_name, req.session.customerId, {
                    maxAge: 3600000
                }).send('Cookie is set');
                console.log("Sent cookies.");
                res.end();
                return;
            });
    });

    module.exports = app;
}());
