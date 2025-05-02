FROM golang:1.24-alpine

RUN go version 

ENV GOPATH=/

COPY ./ ./ 

RUN go mod download
RUN go build -o sso-auth ./cmd/main.go


CMD [ "./sso-auth" ]