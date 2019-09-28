# ############################
# # STEP 1 build executable binary
# ############################
FROM golang:alpine AS builder
# Install git.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/kucar/coinlib
COPY . .
COPY config/config.json /go/bin/gocoin/config/config.json
# Fetch dependencies using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /go/bin/gocoin/main examples/main.go

# # ############################
# # # STEP 2 build a small image
# # ############################
FROM ubuntu:latest
RUN  apt update && apt install -y git
COPY --from=builder /go/bin/gocoin/main /go/bin/gocoin/main
COPY --from=builder /go/bin/gocoin/config/config.json /go/bin/gocoin/config/config.json
EXPOSE 9090
WORKDIR /go/bin/gocoin/
ENTRYPOINT ["/go/bin/gocoin/main"]


#container run : docker run -it -p 9090:9090 name:id 
