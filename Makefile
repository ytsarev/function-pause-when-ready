.PHONY: test
test:
	go test -v -cover .

.PHONY: render
render:
	crossplane beta render example/xr.yaml example/composition.yaml example/functions.yaml -o example/observed.yaml -r

.PHONY: lint
lint:
	golangci-lint run
