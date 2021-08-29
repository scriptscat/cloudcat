GOOS=linux
GOARCH=amd64
VERSION=v0.1.0
DOCKER_REPO=codfrm
REMOTE_REPO=$(DOCKER_REPO)/cloudcat:$(VERSION)

test:
	go test -v ./...

build:
	CGO_LDFLAGS="-static" go build -o cloudcat ./cmd/app

target:
	CGO_LDFLAGS="-static" go build -o cloudcat-$(VERSION)-$(GOOS)-$(GOARCH)/cloudcat ./cmd/app

docker:
	docker build -t cloudcat .

docker-test:
	docker run -it -v $(PWD)/bilibili.zip:/cloudcat/bilibili.zip -v /etc/localtime:/etc/localtime -v /etc/timezone:/etc/timezone $(REMOTE_REPO) exec bilibili.zip

docker-push: docker
	docker tag cloudcat $(REMOTE_REPO)
	docker push $(REMOTE_REPO)

swagger:
	swag init -g internal/interface/http/apiv1/router.go
