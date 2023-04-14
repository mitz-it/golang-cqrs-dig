.PHONY: test


test:
	rm -rf out && \
	mkdir -p out && \
	go test -race -covermode=atomic -coverprofile=out/coverage.out $(go list ./... | grep -v vendor/)