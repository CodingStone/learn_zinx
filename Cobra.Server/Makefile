.PHONY: build #声明build是"伪目标"之后，make就不会去检查是否存在一个叫做build的文件，而是每次运行都执行对应的命令。

#Makefile中变量定义有两种方式[= 与 :=]  := 必须是先定义后使用。  = 可以先使用后定义
SERVICE := zinx
CUR_PWD := $(shell pwd)

SERVER_DEMO_PATH := $(CUR_PWD)/examples/cobra.server
SERVER_DEMO_BIN := $(CUR_PWD)/bin


AUTHOR := Rock
VERSION := 0.0.1
BUILD_INFO := pre
BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M:%S)

#向下传递环境变量
export GO111MODULE=on

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

default: build

build:
	go build  -ldflags $(LD_FLAGS) -gcflags "-N"  -o $(SERVER_DEMO_BIN) $(SERVER_DEMO_PATH)/main.go
clean:
	rm $(SERVER_DEMO_BIN)