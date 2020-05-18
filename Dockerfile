FROM golang:1.14-alpine AS build-env

WORKDIR /var/app
COPY . /var/app

RUN CGO_ENABLED=0 \
	go build \
	-ldflags "-X main.tag=$TAG" \
	-o oauth-revokerd .

###

FROM alpine

# We need ca-certificates to securely retrieve JWKSes
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=build-env /var/app/oauth-revokerd .

EXPOSE 8080 3320 3322
CMD ./oauth-revokerd