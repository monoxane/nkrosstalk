FROM golang:latest

WORKDIR /nkrosstalk
ADD . .
RUN go build src/main.go

EXPOSE 9999
CMD ./main