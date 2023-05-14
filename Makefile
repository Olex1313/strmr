EXEC_NAME=strmr

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build a executable binary
build:
	go build -o $(EXEC_NAME) cmd/main.go

test: mocks build
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

## fresh: Clean everything and launch app
fresh: clean build run

## run: Run app
run: build
	./$(EXEC_NAME)

## clean: Clean everything
clean:
	go clean
	rm -f $(EXEC_NAME)
	rm -rf $(MOCKS_DESTINATION)

MOCKS_DESTINATION=mocks

.PHONY: test build
