BINARY=go-tuimg
VERSION=1.0.0
BUILD_DIR    := ./output

# 构建前自动创建输出目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: all clean linux macos windows

all: linux macos windows

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-macos-arm64 .

macos-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-macos-amd64 .

macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-macos-arm64 .

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .

local:
	go build -o $(BUILD_DIR)/$(BINARY) .

clean:
	rm -rf $(BUILD_DIR)

help:
	@echo "用法:"
	@echo "  make local         - 编译当前系统版本"
	@echo "  make linux         - 编译Linux amd64 + arm64"
	@echo "  make macos         - 编译macOS amd64 + arm64"
	@echo "  make windows       - 编译Windows amd64"
	@echo "  make clean         - 清理编译文件"