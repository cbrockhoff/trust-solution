# FROM weaveworksdemos/msd-java:jre-latest
# FROM openjdk:8u212-jre
FROM adoptopenjdk:12.0.1_12-jre-hotspot

ARG SERVICE_USER='the_user'
ARG SERVICE_GROUP='the_group'
ARG SERVICE_UID='10001'
ARG SERVICE_GID='10001'

RUN groupadd -g ${SERVICE_GID} ${SERVICE_GROUP}
RUN useradd -g ${SERVICE_GID} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

# RUN mkdir /lib64
RUN ln -s /opt/java/openjdk/lib/libjli.so /lib/libjvm.so
# RUN ln -s /usr/lib/jvm/java-1.8-openjdk/lib/amd64/jli/libjli.so /lib/libjli.so
RUN apt-get update
RUN apt-get -y install libcap2-bin
RUN setcap "cap_net_bind_service=+ep" $(readlink -f $(which java))

RUN find / -name "cacerts"
RUN apt-get -y install ca-certificates

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

RUN keytool -importcert -file /usr/local/share/ca-certificates/carts.crt -alias carts -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt
RUN keytool -importcert -file /usr/local/share/ca-certificates/catalogue.crt -alias catalogue -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt
RUN keytool -importcert -file /usr/local/share/ca-certificates/orders.crt -alias orders -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt
RUN keytool -importcert -file /usr/local/share/ca-certificates/payment.crt -alias payment -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt
RUN keytool -importcert -file /usr/local/share/ca-certificates/shipping.crt -alias shipping -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt
RUN keytool -importcert -file /usr/local/share/ca-certificates/user.crt -alias user -keystore /opt/java/openjdk/lib/security/cacerts -storepass changeit -noprompt

WORKDIR /usr/src/app
COPY *.jar ./app.jar
RUN	chown -R ${SERVICE_USER}:${SERVICE_GROUP} ./app.jar

USER ${SERVICE_USER}
ENTRYPOINT ["java", "-Xms64m", "-Xmx128m", "-XX:+UseG1GC", "-Djava.security.egd=file:/dev/urandom", "-Dspring.zipkin.enabled=false", "-jar", "./app.jar", "--port=80"]
