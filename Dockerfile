FROM golang:1.13.4

RUN mkdir /csv-webapp; mkdir /csv-webapp/files; mkdir /csv-webapp/logs

ADD . /csv-webapp

WORKDIR /csv-webapp

RUN ./build.sh

CMD ["/csv-webapp/bin/./webapp"]