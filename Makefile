BINARY_NAME=app
BUILD_DIR=bin

build:
	mkdir -p ${BUILD_DIR}
	cd app && go build -ldflags="-s -w" -o ../${BUILD_DIR}/${BINARY_NAME}