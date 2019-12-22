FROM golang:1.13.4

RUN mkdir /csv-webapp; mkdir /csv-webapp/files; mkdir /csv-webapp/logs

ADD . /csv-webapp

COPY ./scripts/a-mongo.js /docker-entrypoint-initdb.d/
COPY ./scripts/b-mongo.js /docker-entrypoint-initdb.d/

WORKDIR /csv-webapp

RUN ./build.sh

CMD ["/csv-webapp/bin/./webapp"]