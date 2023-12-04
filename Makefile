.PHONY: test
test:
	crossplane beta render example/xr.yaml example/composition.yaml example/functions.yaml -o example/observed.yaml -r
