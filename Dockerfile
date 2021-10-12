FROM golang:alpine AS build
# run build process in /app directory 
WORKDIR /app
# copy dependencies and get them
COPY ./go.mod ./go.sum ./
RUN go get -d -v ./...
# copy go src file(s)
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
# build the binary
RUN GOOS=linux GOARCH=amd64 go build -o kvs ./cmd/main.go

FROM alpine
# copy the binary from build stage
COPY --from=build /app/kvs /bin/kvs
# use non root
USER 1000:1000
# start server
CMD ["./bin/kvs"]