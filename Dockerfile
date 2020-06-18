FROM golang:1.14-alpine3.12 AS build

RUN apk --update add \
		ca-certificates \
		gcc \
		git \
		musl-dev

RUN echo 'nobody:x:65534:65534:nobody:/:' > /tmp/passwd \
	&& echo 'nobody:x:65534:' > /tmp/group

COPY go.mod go.sum /go/src/github.com/juli3nk/faas-connector-stan/
WORKDIR /go/src/github.com/juli3nk/faas-connector-stan

ENV GO111MODULE on
RUN go mod download

COPY stan stan
COPY config config
COPY main.go .

# Stripping via -ldflags "-s -w"
#RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /usr/bin/producer
RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o /tmp/producer


FROM scratch

COPY --from=build /tmp/group /tmp/passwd /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /tmp/producer /producer

USER nobody:nobody

ENTRYPOINT ["/producer"]
