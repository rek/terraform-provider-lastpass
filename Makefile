GO=/usr/local/go/bin/go
TEST?=$$(${GO} list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=rezroo
NAME=lastpass
BINARY=terraform-provider-${NAME}
VERSION=0.6.0
OS_ARCH=linux_amd64

default: install

build:
	${GO} build -o ${BINARY}

release:
	# GOOS=darwin GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	# GOOS=darwin GOARCH=arm64 ${GO} build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	# GOOS=freebsd GOARCH=386 ${GO} build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	# GOOS=freebsd GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	# GOOS=freebsd GOARCH=arm ${GO} build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	# GOOS=linux GOARCH=386 ${GO} build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	# GOOS=linux GOARCH=arm ${GO} build -o ./bin/${BINARY}_${VERSION}_linux_arm
	# GOOS=openbsd GOARCH=386 ${GO} build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	# GOOS=openbsd GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	# GOOS=solaris GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	# GOOS=windows GOARCH=386 ${GO} build -o ./bin/${BINARY}_${VERSION}_windows_386
	# GOOS=windows GOARCH=amd64 ${GO} build -o ./bin/${BINARY}_${VERSION}_windows_amd64

# also cp to root plugins folder to make backward compatible with TF 0.12
install: build ~/.terraform.d/plugins/${BINARY}_v${VERSION}
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

~/.terraform.d/plugins/${BINARY}_v${VERSION}: ${BINARY}
	cp ${BINARY} ~/.terraform.d/plugins/${BINARY}_v${VERSION}

test: 
	${GO} test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: 
	TF_ACC=1 ${GO} test $(TEST) -v $(TESTARGS) -timeout 120m
