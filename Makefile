help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build a executable binary
build:
	go build cmd/main.go

test: mocks build
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

## fresh: Clean everything and launch app
fresh: clean build run

## run: Run app
run: build
	./main

## clean: Clean everything
clean:
	go clean
	rm -f **/wire_gen.go
	rm -f main
	rm -rf $(MOCKS_DESTINATION)

MOCKS_DESTINATION=mocks

.PHONY: test build
