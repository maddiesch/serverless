.PHONY: test
test:
	go test -v ./...

.PHONY: dynamodb
dynamodb:
	docker run -p 8000:8000 amazon/dynamodb-local:latest
