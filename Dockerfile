FROM golang:1.16 as build

ARG GOOS=linux
ARG GOARCH=amd64
ARG GOPROXY="https://goproxy.cn,direct"

WORKDIR /cloudcat

COPY . .

# TODO: 交叉编译

RUN make build

ARG ARCH

FROM ${ARCH}busybox:1.33.1-musl

WORKDIR /cloudcat

COPY --from=build /cloudcat/cloudcat .

RUN ls -l && chmod +x cloudcat

ENTRYPOINT ["./cloudcat"]
