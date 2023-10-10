NAME=sendoracity
BINARY=terraform-provider-${NAME}
HOSTNAME=registry.sendora.io
NAMESPACE=pchanvallon
VERSION=0.0.1
BASE_URI=http://localhost:8080

ifndef GOOS
	GOOS=linux
endif
ifndef GOARCH
	GOARCH=amd64
endif

default: install

tidy:
	go mod tidy

build: tidy
	go build -o ${BINARY}

build_debug: tidy
	go build -gcflags=all="-N -l" -o ${BINARY}

mv:
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/${BINARY}

install: build mv

debug: build_debug mv

.PHONY: testacc
testacc:
	BASE_URI=${BASE_URI} TF_ACC=1 go test ./...
