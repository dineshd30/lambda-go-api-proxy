# Default target
.PHONY: all
all: build

# build go service
.PHONY: build
build:
	docker build --platform linux/arm64 -t gombroxy -f Dockerfile . 

# run go service
.PHONY: run
run: build
	docker run -p 9000:8080 gombroxy


