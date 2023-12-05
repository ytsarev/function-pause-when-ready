REPO ?= xnull/function-pause-when-ready
TAG ?= v0.1.2

.PHONY: test
test:
	go test -v -cover .

.PHONY: render
render:
	crossplane beta render example/xr.yaml example/composition.yaml example/functions.yaml -o example/observed.yaml -r

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	docker build . --quiet --platform=linux/arm64 --tag runtime-arm64
	crossplane xpkg build --package-root=package --embed-runtime-image=runtime-arm64 --package-file=function-arm64.xpkg
	docker build . --quiet --platform=linux/amd64 --tag runtime-amd64
	crossplane xpkg build --package-root=package --embed-runtime-image=runtime-amd64 --package-file=function-amd64.xpkg

.PHONY: push
push:
	crossplane xpkg push --package-files=function-amd64.xpkg,function-arm64.xpkg $(REPO):$(TAG)
