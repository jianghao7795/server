.PHONY: build dev start clean tool help

all: build

build:
	 GOEXPERIMENT=jsonv2 CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -trimpath -o bin/fiber -ldflags="-s -w" cmd/main.go 

start:
	air

dev:
	go run cmd/main.go

local:
	go run . -c ./conf/

tool:
	go tool vet . |& grep -v vendor; true
	gofmt -w .

clean:
	go clean -i .

air:
	air -c .air.toml

test:
	go test -v ./...

wire:
	cd core && wire

# === Podman 部署 ===
pod-create:
	-podman pod create --name man-pod -p 3100:3100 -p 3306:3306 -p 6379:6379

image:
	podman build -t man-app .

up: pod-create
	podman run -d --pod man-pod --name man-mysql \
		-e MYSQL_ROOT_PASSWORD=123456 \
		-e MYSQL_DATABASE=server-fiber \
		-v mysql-data:/var/lib/mysql \
		docker.io/library/mysql:8.0
	podman run -d --pod man-pod --name man-redis \
		docker.io/library/redis:7-alpine
	@echo "等待 MySQL 就绪..."
	@sleep 10
	podman run -d --pod man-pod --name man-app man-app
	@echo "全部启动完成: http://localhost:3100"

down:
	-podman stop man-app man-mysql man-redis
	-podman rm man-app man-mysql man-redis
	-podman pod rm man-pod

pod-clean: down
	-podman volume rm mysql-data

help:
	@echo "make build: 本地编译"
	@echo "make dev: go run ."
	@echo "make image: 构建 podman 镜像"
	@echo "make up: 创建 pod 并启动 mysql/redis/app"
	@echo "make down: 停止并删除所有容器和 pod"
	@echo "make pod-clean: 完全清理（含数据卷）"
