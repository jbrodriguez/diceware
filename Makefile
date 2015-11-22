#
# Makefile to perform "live code reloading" after changes to .go files.
#
# n.b. you must install fswatch (OS X: `brew install fswatch`)
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)

# binary name to kill/restart
PROG = diceware
 
# targets not associated with files
.PHONY: dependencies default build test coverage clean kill restart serve
 
# check we have a couple of dependencies
dependencies:
	@command -v fswatch --version >/dev/null 2>&1 || { printf >&2 "fswatch is not installed, please run: brew install fswatch\n"; exit 1; }
 
# default targets to run when only running `make`
default: dependencies test
 
# clean up
clean:
	go clean
 
# run formatting tool and build
build: dependencies clean
	go build fmt
	go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v

# run unit tests with code coverage
test: dependencies
	go test -v
 
# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage.out
 
# attempt to kill running server
kill:
	-@killall -9 $(PROG) 2>/dev/null || true
 
# @make build; (if [ "$$?" -eq 0 ]; then (env GIN_MODE=debug ./${PROG} -pg_port_5432_tcp_addr api.apertoire.org &); fi)
# attempt to build and start server
restart:
	@make kill
	@make build; (if [ "$$?" -eq 0 ]; then (env GIN_MODE=debug ./${PROG} -pg_port_5432_tcp_addr localhost &); fi)
 
# watch .go files for changes then recompile & try to start server
# will also kill server after ctrl+c
serve: dependencies
	@make restart 
	@fswatch -o ./*.go | xargs -n1 -I{} make restart || make kill