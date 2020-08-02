BUILD_DIR=bin/
EXENAME=urlshort

GO=$(shell which go)

# Check OS 
UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
       GOOS=darwin
else
    ifeq ($(UNAME), Linux)
        GOOS=linux
    endif
endif

.PHONY: test run build clean

test:
	$(GO) test -v ./pkg/...

run: test
	$(GO) run -race cmd/main.go -test

# Build the project
build: test
	 CGO_ENABLED=0 GOOS=$(GOOS) GO111MODULES=auto $(GO) build -a -installsuffix cgo -o \
				 $(BUILD_DIR)$(EXENAME) cmd/main.go

# Clean out the build dir
clean: $(BUILD_DIR)
	 [ -d $< ] && rm -r $< && $(GO) clean

# Make build directory if it doesn't already exist
$(BUILD_DIR):
	[ -d $@ ] || mkdir -p $@
