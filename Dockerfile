FROM golang:1.14-alpine AS build-env

WORKDIR /var/app
COPY . /var/app

RUN go build -o ./oauth-revokerd -ldflags "-X main.tag=$TAG" .

###

FROM alpine

# We need ca-certificates to securely retrieve JWKSes
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=build-env /var/app/oauth-revokerd .

EXPOSE 8080
CMD ./oauth-revokerd