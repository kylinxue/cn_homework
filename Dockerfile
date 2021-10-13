FROM ubuntu
ENV MY_SERVICE_PORT=80
LABEL server.name=my_http_server
COPY bin/amd64/httpserver /httpserver
EXPOSE ${MY_SERVICE_PORT}
ENTRYPOINT /httpserver