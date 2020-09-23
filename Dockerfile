FROM ubuntu:latest

COPY ./myNode /usr/local/bin/

ENTRYPOINT ["myNode"]
