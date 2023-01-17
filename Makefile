.PHONY: test
test:
	go test -v -failfast -race -coverpkg=./... -covermode=atomic

.PHONY: bench
bench:
	go test -v -failfast -bench=. -benchmem -run=^$$
