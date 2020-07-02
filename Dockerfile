FROM golang:1-alpine as builder

RUN apk --update upgrade \
&& apk --no-cache --no-progress add git make \
&& rm -rf /var/cache/apk/*

WORKDIR /go/kutteri

ENV GO111MODULE on

# Download go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.6
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates git \
    && rm -rf /var/cache/apk/*

COPY --from=builder /go/kutteri/kutteri /usr/bin/kutteri

ENTRYPOINT ["/usr/bin/kutteri"]
