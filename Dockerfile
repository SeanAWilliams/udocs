FROM golang:1.6-alpine

RUN apk update && apk add --no-cache --virtual .build-deps bash git openssh

ENV GOPATH $HOME/go
ENV PATH $PATH:$GOPATH/bin

RUN mkdir -p $GOPATH/src/github.com/ultimatesoftware/udocs
WORKDIR $GOPATH/src/github.com/ultimatesoftware/udocs

COPY ./ ./

RUN ./bin/install.sh && udocs env

EXPOSE 9554
CMD ["udocs", "serve", "--headless"]
