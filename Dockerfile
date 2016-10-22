FROM golang:1.7-alpine

RUN apk update && apk add --no-cache --virtual .build-deps bash git openssh gcc

ENV GOPATH $HOME/go
ENV PATH $PATH:$GOPATH/bin

RUN mkdir -p $GOPATH/src/github.com/ultimatesoftware/udocs
WORKDIR $GOPATH/src/github.com/ultimatesoftware/udocs
RUN go env
COPY ./ ./
RUN rm -rf $GOPATH/src/github.com/ultimatesoftware/udocs/vendor
RUN go get -v ./... && go build && go install
EXPOSE 9554
CMD ["udocs", "serve", "--headless"]
