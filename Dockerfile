FROM golang:latest

RUN mkdir -p $GOPATH/src/github.com/wasuken/compiew
WORKDIR $GOPATH/src/github.com/wasuken/compiew
COPY . $GOPATH/src/github.com/wasuken/compiew

RUN go mod tidy
