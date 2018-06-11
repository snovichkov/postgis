.PHONY: all install-deps update-deps test test-with-coverage test-with-coverage-profile lint lint-format lint-import lint-style

# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_FILE_NAME ?= coverage.out

# Set a default `min_confidence` value for `golint`
GO_LINT_MIN_CONFIDENCE ?= 0.2

all: install-deps test

install-deps:
	@echo "Install dependencies"
	@go get golang.org/x/lint/golint
	@go get github.com/stretchr/testify
	@go get golang.org/x/tools/cmd/goimports

update-deps:
	@echo "Update dependencies"
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/stretchr/testify
	@go get -u golang.org/x/tools/cmd/goimports

test:
	@echo "Run unit tests"
	@go test -v ./...

test-with-coverage:
	@echo "Run unit tests with coverage"
	@go test -cover ./...

test-with-coverage-profile:
	@echo "Run unit tests with coverage profile"
	@ERR=0; \
	echo "mode: ${GO_TEST_COVERAGE_MODE}" > "${GO_TEST_COVERAGE_FILE_NAME}"; \
	for package in $$(go list ./...); do \
		go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile "coverage_$${package##*/}.out" "$${package}" || { \
        	ERR=$$?; \
          	break; \
        }; \
        if [ ! -f "coverage_$${package##*/}.out" ]; then \
        	continue; \
		fi; \
		{ \
			sed '1d' "coverage_$${package##*/}.out" >> "${GO_TEST_COVERAGE_FILE_NAME}" && \
			rm "coverage_$${package##*/}.out"; \
		} || { \
			ERR=$$?; \
			break; \
		}; \
	done; \
	if [ $$ERR != 0 ]; then \
		exit $$ERR; \
	fi;
	@echo "Generate coverage report"
	@go tool cover -func="${GO_TEST_COVERAGE_FILE_NAME}"; \
	rm "${GO_TEST_COVERAGE_FILE_NAME}";

lint: lint-format lint-import

lint-format:
	@echo "Check formatting"
	@errors=$$(gofmt -l ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

lint-import:
	@echo "Check imports"
	@errors=$$(goimports -l $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

lint-style:
	@echo "Check code style"
	@errors=$$(golint -min_confidence=${GO_LINT_MIN_CONFIDENCE} $$(go list ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi
