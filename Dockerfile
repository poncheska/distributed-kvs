FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bin/cluster ./cmd/cluster/

FROM alpine

COPY --from=builder /app/bin/cluster /cluster

CMD [ "/cluster" ]
