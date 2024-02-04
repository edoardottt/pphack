REPO=github.com/edoardottt/pphack

remod:
	@rm -rf go.*
	@go mod init ${REPO}
	@go get ./...
	@go mod tidy -v
	@echo "Done."

update:
	@go get -u ./...
	@go mod tidy -v
	@echo "Done."

lint:
	@golangci-lint run

build:
	@go build ./cmd/pphack/
	@sudo mv pphack /usr/local/bin/
	@echo "Done."

clean:
	@sudo rm -rf /usr/local/bin/pphack
	@echo "Done."

test:
	@go test -race ./...
	@echo "Done."