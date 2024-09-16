FROM golang:1.22-alpine AS builder
RUN apk add --update --no-cache git
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -mod=mod -a -o /aws-azrebalance-controller cmd/aws-azrebalance-controller/main.go

FROM scratch
COPY --from=builder /aws-azrebalance-controller /aws-azrebalance-controller
CMD ["./aws-azrebalance-controller"]