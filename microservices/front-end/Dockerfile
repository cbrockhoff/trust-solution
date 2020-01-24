FROM node:10-alpine
ENV NODE_ENV "production"
ENV PORT 8079
EXPOSE 8079
RUN addgroup mygroup && adduser -D -G mygroup myuser && mkdir -p /usr/src/app && chown -R myuser /usr/src/app
RUN apk update && apk add ca-certificates

COPY certs/carts.crt /usr/local/share/ca-certificates/carts.crt
COPY certs/catalogue.crt /usr/local/share/ca-certificates/catalogue.crt
COPY certs/orders.crt /usr/local/share/ca-certificates/orders.crt
COPY certs/payment.crt /usr/local/share/ca-certificates/payment.crt
COPY certs/shipping.crt /usr/local/share/ca-certificates/shipping.crt
COPY certs/user.crt /usr/local/share/ca-certificates/user.crt

RUN chmod 644 /usr/local/share/ca-certificates/carts.crt
RUN chmod 644 /usr/local/share/ca-certificates/catalogue.crt
RUN chmod 644 /usr/local/share/ca-certificates/orders.crt
RUN chmod 644 /usr/local/share/ca-certificates/payment.crt
RUN chmod 644 /usr/local/share/ca-certificates/shipping.crt
RUN chmod 644 /usr/local/share/ca-certificates/user.crt

RUN update-ca-certificates

# Prepare app directory
WORKDIR /usr/src/app
COPY package.json /usr/src/app/
COPY yarn.lock /usr/src/app/
RUN chown myuser /usr/src/app/yarn.lock
RUN apk --no-cache add curl

USER myuser
RUN yarn install

COPY . /usr/src/app

ENV NODE_TLS_REJECT_UNAUTHORIZED 0

# Start the app
CMD ["/usr/local/bin/npm", "start"]
