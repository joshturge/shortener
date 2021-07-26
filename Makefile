BUILD_DIR=bin/
EXENAME=shortener

GO=`which go`

.PHONY: test run build clean

test:
	$(GO) test -v ./pkg/...

run: test
	$(GO) run cmd/main.go -test

# Build the project
build: test
	 CGO_ENABLED=0 GO111MODULES=auto $(GO) build -a -installsuffix cgo -o \
				 $(BUILD_DIR)$(EXENAME) cmd/main.go

# Clean out the build dir
clean:
	rm -rf $(BUILD_DIR) && $(GO) clean
