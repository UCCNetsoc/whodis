FROM golang:latest AS dev

WORKDIR /bot

RUN apk add git

RUN GO111MODULE=on go get github.om/cortesi/modd/cmd/modd

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install github.com/uccnetsoc/veribot

RUN ["go", "run", "*.go"]


FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/veribot ./veribot

CMD ["sh", "-c", "veribot"]