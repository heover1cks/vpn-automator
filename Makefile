APPLICATION=vpn-automator
GO      ?= go
GOARCH	?= $(shell $(GO) env GOARCH)
GOOS	?= $(shell $(GO) env GOOS)

build:
	GOARCH=${GOARCH} GOOS=${GOOS} go build -o ${APPLICATION} cmd/vpn-automator.go
	cp vpn-automator /usr/local/bin