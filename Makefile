F=

# run tests
.PHONY: t
t:
	@go test ./... -run "${F}" \
		| grep -v "no tests to run" \
		| grep -v "no test files"

# run all tests & generate coverage
.PHONY: c
c:
	@go test -count=1 -race -covermode=atomic ./... -coverprofile=cover.out > /dev/null && \
	go tool cover -func cover.out \
		| grep -v '[89]\d\.\d%' | grep -v '100.0%' \
		| grep -v 'log/noop.go' \
		| grep -v 'log/kv_parser.go' \
		|| true
	@go tool cover -html=cover.out
	@rm cover.out
