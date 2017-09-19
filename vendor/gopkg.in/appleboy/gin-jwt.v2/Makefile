.PHONY: test

test:
	go test -v -covermode=count -coverprofile=coverage.txt

html:
	go tool cover -html=coverage.txt
