FROM golang:1.23.2-alpine3.20
LABEL maintainer="LasramR <pro@remi-marsal.com>"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/sclng-backend-test-lasramR

CMD $GOPATH/bin/sclng-backend-test-lasramR
