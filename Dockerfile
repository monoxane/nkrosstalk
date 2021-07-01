FROM golang:latest

WORKDIR /nkrosstalk
ADD . .
RUN go build -o nkrosstalk src/main.go src/nk.go

EXPOSE 7788
CMD ./nkrosstalk