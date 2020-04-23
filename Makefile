all: deps unitTest e2eTest
test: unitTest e2eTest

deps:
	go get -d -v github.com/google/uuid
	go get -d -v github.com/golang/mock/gomock
	go get -d -v github.com/onsi/ginkgo
	go get -d -v github.com/onsi/gomega

unitTest:
	go test -v ./... -tags=unit

e2eTest:
	go test -v ./... -tags=e2e

.PHONY: deps unitTest e2eTest
