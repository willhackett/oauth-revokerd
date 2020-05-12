include .env

export VERSION := $(or $(CIRCLE_TAG),$(shell git log --pretty=format:'%h' -n 1))
export $(shell sed 's/=.*//' .env)

GO_PACKAGES = . ./app/...

build:
	go build -o ./oauth-revokerd -ldflags "-X main.version=$(VERSION)" .

test:
	golint -set_exit_status ${GO_PACKAGES}
	go vet ${GO_PACKAGES}
	go test ${GO_PACKAGES}

