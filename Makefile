GOOS=linux
GOARCH=amd64
VERSION=v0.1.0
DOCKER_REPO=codfrm
REMOTE_REPO=$(DOCKER_REPO)/cloudcat:$(VERSION)

NAME=cloudcat-$(VERSION)-$(GOOS)-$(GOARCH)/cloudcat
SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

swagger:
	swag fmt -g interfaces/api/router.go
	swag init -g interfaces/api/router.go

test:
	GOOS=$(GOOS) go test -v ./...

generate:
	go generate ./... -x

build: swagger generate
	CGO_LDFLAGS="-static" go build -tags netgo -o cloudcat$(SUFFIX) ./cmd/cloudcat

target:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_LDFLAGS="-static" go build -tags netgo -o $(NAME)$(SUFFIX) ./cmd/cloudcat

docker:
	docker build -t cloudcat .

docker-test:
	docker run -it -v $(PWD)/bilibili.zip:/cloudcat/bilibili.zip -v /etc/localtime:/etc/localtime -v /etc/timezone:/etc/timezone $(REMOTE_REPO) exec bilibili.zip

docker-push: docker
	docker tag cloudcat $(REMOTE_REPO)
	docker push $(REMOTE_REPO)
