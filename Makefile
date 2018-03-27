build:
	go install -v

depends:
	dep ensure

dev-depends:
	go get -u github.com/smartystreets/goconvey

test:
	goconvey