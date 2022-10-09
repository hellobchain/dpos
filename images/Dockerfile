FROM golang:1.16 as builder
ARG Version
ARG CommitVersion
ARG BuildTime
LABEL version=$Version comshbuimit=$CommitVersion create_time=$BuildTime

ADD . /dpos
WORKDIR /dpos
RUN  go version && go env && gcc -v && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build \
     --mod=vendor  -v -o dpos cmd/main.go
#生成中间镜像后,将build之后的可执行文件考到新的镜像中
FROM alpine:3.14 as dpos
ARG Version
ARG CommitVersion
ARG BuildTime
LABEL version=$Version commit=$CommitVersion create_time=$BuildTime
COPY --from=builder  /dpos/dpos /usr/local/bin
# 切换软件源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /dpos/dpos
#容器内部开放端口
CMD ["dpos"]
