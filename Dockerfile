FROM golang:1.14-alpine as build
LABEL maintainer=""
ENV GO111MODULE=on
WORKDIR /harukax

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/harukax

FROM alpine:3.9

COPY --from=build /build/harukax /harukax
RUN touch .env # godotenv is weird and needs this

ENTRYPOINT ["/harukax"]
