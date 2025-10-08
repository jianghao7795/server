FROM golang:latest AS builder

LABEL org.opencontainers.image.authors="jianghao"

ENV GOPROXY=https://goproxy.cn,direct
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
ENV GO111MODULE=on

WORKDIR /app
COPY . /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -trimpath -o fiber -ldflags="-s -w" cmd/main.go #CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o fiber cmd/main.go

FROM rockylinux:9-minimal AS runner
WORKDIR /app

COPY --from=builder /app/fiber .
COPY --from=builder /app/config.yaml ./conf/config.yaml
COPY --from=builder /app/rbac_model.conf .
COPY --from=builder /app/docs/ ./docs/

EXPOSE 3100
CMD ["/app/fiber", "-c", "./conf/"]
# ENTRYPOINT ["/app/server-fiber", "-c", "config.yaml"]


#
# 1. 关键参数说明
# CGO_ENABLED=0
# 禁用 CGO，强制生成纯静态二进制文件（不依赖动态链接库），适用于无 C 代码依赖的场景。
#
# GOOS=linux
# 指定目标操作系统为 Linux，支持跨平台编译。
#
# -tags=jsoniter
# 启用 jsoniter 标签，使用 jsoniter 库替代 Go 标准库的 encoding/json/v2，适用于需要高性能 JSON 解析的场景。
#
# -a
# 强制重新编译所有依赖包（即使已是最新），避免因缓存导致的潜在问题，但会增加编译时间。
#
# -installsuffix cgo
# 已弃用。此参数原用于隔离不同编译选项的包缓存，但在 Go 1.10+ 版本后，-trimpath 和模块机制已替代其功能，可移除。
#
# -ldflags "-w"
# 禁用 DWARF 调试信息生成，减少二进制文件体积（但未移除符号表）。
