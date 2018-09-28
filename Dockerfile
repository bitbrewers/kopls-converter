FROM golang:1.11 as builder

WORKDIR /go/src/github.com/bitbrewers/kopls-converter

RUN go get github.com/golang/dep/cmd/dep
ADD Gopkg.toml Gopkg.lock ./
RUN dep ensure -v --vendor-only

COPY server server
COPY templates templates
COPY Makefile *.go ./
RUN make build

FROM scratch
COPY migrations migrations
COPY --from=builder /go/src/github.com/bitbrewers/kopls-converter/builds/kopls-converter /kopls-converter

ENTRYPOINT [ "/kopls-converter" ]
