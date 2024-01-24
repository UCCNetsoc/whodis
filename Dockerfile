FROM golang:1.21.6-alpine AS dev

WORKDIR /app

RUN apk add git

# RUN GO111MODULE=on go get github.com/cortesi/modd/cmd/modd

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Get latest tag
RUN git describe --abbrev=0 > /version

RUN go install github.com/UCCNetsoc/whodis/cmd/whodis

CMD [ "go", "run", "*.go" ]

FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/whodis ./whodis
COPY --from=dev /version /version

CMD ["sh", "-c", "export BOT_VERSION=$(cat /version) && whodis -p"]
