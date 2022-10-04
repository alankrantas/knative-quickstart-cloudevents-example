FROM golang:alpine as builder
WORKDIR /app

COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -mod=readonly -v -o server

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server

EXPOSE 8080
CMD ./server