PROJECTNAME1=$(shell basename "$(PWD)")
PROJECTNAME=$(PROJECTNAME1:go-%=%)
APPNAME=$(patsubst "%",%,$(shell grep -E "AppName[ \t]+=[ \t]+" doc.go|grep -Eo "\\\".+\\\""))
VERSION=$(shell grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+")
include .env
-include .env.local
# ref: https://kodfabrik.com/journal/a-good-makefile-for-go/

  # https://www.gnu.org/savannah-checkouts/gnu/make/manual/html_node/Text-Functions.html
  # https://stackoverflow.com/questions/19571391/remove-prefix-with-make


# Go related variables.
GOBASE       =  $(shell pwd)
#,#GOPATH="$(GOBASE)/vendor:$(GOBASE)"
#,#GOPATH=$(GOBASE)/vendor:$(GOBASE):$(shell dirname $(GOBASE))
#GOPATH2     =  $(shell dirname $(GOBASE))
#GOPATH1     =  $(shell dirname $(GOPATH2))
#GOPATH0     =  $(shell dirname $(GOPATH1))
#GOPATH      =  $(shell dirname $(GOPATH0))
GOBIN        =  $(GOBASE)/bin
GOFILES      =  $(wildcard *.go)
SRCS         =  $(shell git ls-files '*.go')
PKGS         =  $(shell go list ./...)
GIT_VERSION  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_REVISION := $(shell git rev-parse --short HEAD)
#GITHASH     =  $(shell git rev-parse HEAD)
#BUILDTIME   := $(shell date "+%Y%m%d_%H%M%S")
#BUILDTIME   =  $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILDTIME    =  $(shell date -u '+%Y-%m-%d_%H-%M-%S')
GOVERSION    =  $(shell go version)
BIN          =  $(GOPATH)/bin
GOLINT       =  $(BIN)/golint
GOCYCLO      =  $(BIN)/gocyclo
GOYOLO       =  $(BIN)/yolo


# GO111MODULE = on
GOPROXY     = $(or $(GOPROXY_CUSTOM),direct)

# Redirect error output to a file, so we can show it in development mode.
STDERR      = $(or $(STDERR_CUSTOM),/tmp/.$(PROJECTNAME)-stderr.txt)

# PID file will keep the process id of the server
PID         = $(or $(PID_CUSTOM),/tmp/.$(PROJECTNAME).pid)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent


goarch=amd64
W_PKG=github.com/hedzr/cmdr/conf
LDFLAGS := -s -w \
	-X '$(W_PKG).Buildstamp=$(BUILDTIME)' \
	-X '$(W_PKG).Githash=$(GIT_REVISION)' \
	-X '$(W_PKG).GoVersion=$(GOVERSION)' \
	-X '$(W_PKG).Version=$(VERSION)'
# -X '$(W_PKG).AppName=$(APPNAME)'
GO := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" \
	GO111MODULE=on GOPROXY=$(GOPROXY) go
GO_OFF := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" \
	GO111MODULE=off go




#
#LDFLAGS=
M = $(shell printf "\033[34;1m▶\033[0m")
ADDR = ":5q5q"
CN = hedzr/$(N)




ifeq ($(OS),Windows_NT)
    LS_OPT=
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    LS_OPT=
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
        LS_OPT=--color
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
        LS_OPT=-G
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif






.PHONY: build compile exec clean
.PHONY: run build-linux build-ci
.PHONY: go-build go-generate go-mod-download go-get go-install go-clean
.PHONY: godoc format fmt lint cov gocov coverage codecov cyclo bench



# For the full list of GOARCH/GOOS, take a look at:
#  https://github.com/golang/go/blob/master/src/go/build/syslist.go
#
# A snapshot is:
#  const goosList = "aix android darwin dragonfly freebsd hurd illumos js linux nacl netbsd openbsd plan9 solaris windows zos "
#  const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm "
#©


## build: Compile the binary. Synonym of `compile`
build: directories compile


## build-win: build to windows executable, for LAN deploy manually.
build-win:
	@-$(MAKE) -s go-build-task os=windows goarchset=amd64

## build-linux: build to linux executable, for LAN deploy manually.
build-linux:
	@-$(MAKE) -s go-build-task os=linux goarchset=amd64

## build-nacl: build to nacl executable, for LAN deploy manually.
build-nacl:
	# NOTE: can't build to nacl with golang 1.14 and darwin
	#    chmod +x $(GOBIN)/$(an)_$(os)_$(goarch)*;
	#    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch)*;
	#    gzip -f $(GOBIN)/$(an)_$(os)_$(goarch);
	#    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch)*;
	@-$(MAKE) -s go-build-task os=nacl goarchset="386 arm amd64p32"
	@echo "  < All Done."
	@ls -la $(LS_OPT) $(GOBIN)/*


## build-plan9: build to plan9 executable, for LAN deploy manually.
build-plan9: goarchset = "386 amd64"
build-plan9:
	@-$(MAKE) -s go-build-task os=plan9 goarchset=$(goarchset)

## build-freebsd: build to freebsd executable, for LAN deploy manually.
build-freebsd:
	@-$(MAKE) -s go-build-task os=freebsd goarchset=amd64

## build-riscv: build to riscv64 executable, for LAN deploy manually.
build-riscv:
	@-$(MAKE) -s go-build-task os=linux goarchset=riscv64

## build-ci: run build-ci task. just for CI tools
build-ci:
	@echo "  >  Building binaries in CI flow..."
	$(foreach os, linux darwin windows, \
	  @-$(MAKE) -s go-build-task os=$(os) goarchset="386 amd64" \
	)
	@echo "  < All Done."
	@ls -la $(LS_OPT) $(GOBIN)/*

go-build-task: directories
	@echo "  >  Building $(os)/$(goarchset) binary..."
	@#echo "  >  LDFLAGS = $(LDFLAGS)"
	# unsupported GOOS/GOARCH pair nacl/386 ??
	$(foreach an, $(MAIN_APPS), \
	  echo "  >  APP NAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)"; \
		$(eval ANAME := $(shell for an1 in $(MAIN_APPS); do \
			if [[ $(an) == $$an1 ]]; then \
			  if [[ $$an1 == cli ]]; then echo $(APPNAME); else echo $$an1; fi; \
			fi; \
		done)) \
	  $(foreach goarch, $(goarchset), \
	    echo "     >> Building (-trimpath) $(GOBIN)/$(ANAME)_$(os)_$(goarch)...$(os)" >/dev/null; \
	    $(GO) build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME)_$(os)_$(goarch) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an)/$(MAIN_ENTRY_FILE); \
	    chmod +x $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	    ls -la $(LS_OPT) $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	) \
	)
	#	$(foreach an, $(MAIN_APPS), \
	#	  $(eval ANAME := $(shell if [ "$(an)" == "cli" ]; then echo $(APPNAME); else echo $(an); fi; )) \
	#	  echo "  >  APP NAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)|ANAME:$(ANAME)"; \
	#	  $(foreach goarch, $(goarchset), \
	#	    echo "     >> Building (-trimpath) $(GOBIN)/$(ANAME)_$(os)_$(goarch)...$(os)" >/dev/null; \
	#	    $(GO) build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME)_$(os)_$(goarch) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an); \
	#	    chmod +x $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	    ls -la $(LS_OPT) $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	) \
	#	)
	#@ls -la $(LS_OPT) $(GOBIN)/*linux*




## compile: Compile the binary.
compile: directories go-clean go-generate
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-build 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/' 1>&2

# @cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2
#@if [[ -z "$(STDERR)" ]]; then echo; else echo -e "\n\nError:\n\n"; cat $(STDERR)  1>&2; fi

## exec: Run given cmd, wrapped with custom GOPATH. eg; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

# go-compile: go-clean go-generate go-build

ooo:
	$(eval ANAME := $(shell for an in $(MAIN_APPS); do \
		if [[ $$an == cli ]]; then A=$(APPNAME); echo $(APPNAME); \
		else A=$$an; echo $$an; \
		fi; \
	done))
	@echo "ANAME = $(ANAME), $$ANAME, $$A"

ox: go-clean go-generate
	$(MAKE) -s go-build

## run: go run xxx
run:
	@$(GO) run -ldflags "$(LDFLAGS)" $(GOBASE)/cli/main.go 

go-build:
	@echo "  >  Building apps: $(MAIN_APPS)..."
	$(foreach an, $(MAIN_APPS), \
		$(eval ANAME := $(shell for an1 in $(MAIN_APPS); do \
			if [[ "$(an)" == $$an1 ]]; then \
			  if [[ $$an1 == cli ]]; then echo $(APPNAME); else echo $$an1; fi; \
			fi; \
		done)) \
	  echo "  >  >  Building $(MAIN_BUILD_PKG)/$(an) -> $(ANAME) ..."; \
	  echo "        +race. -trimpath. APPNAME = $(APPNAME), LDFLAGS = $(LDFLAGS)"; \
	  $(GO) build -v -race -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an)/$(MAIN_ENTRY_FILE); \
	  ls -la $(LS_OPT) $(GOBIN)/$(ANAME); \
	)
	ls -la $(LS_OPT) $(GOBIN)/
	if [[ -d ./plugin/demo ]]; then \
	  $(GO) build -v -race -buildmode=plugin -o ./ci/local/share/fluent/addons/demo.so ./plugin/demo && \
	  chmod +x ./ci/local/share/fluent/addons/demo.so && \
	  ls -la $(LS_OPT) ./ci/local/share/fluent/addons/demo.so; fi
	# go build -o $(GOBIN)/$(APPNAME) $(GOFILES)
	# chmod +x $(GOBIN)/*

go-generate:
	@echo "  >  Generating dependency files ('$(generate)') ..."
	@$(GO) generate $(generate) ./...

go-mod-download:
	@$(GO) mod download

go-get:
	# Runs `go get` internally. e.g; make install get=github.com/foo/bar
	@echo "  >  Checking if there is any missing dependencies...$(get)"
	@$(GO) get $(get)

go-install:
	@$(GO) install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@$(GO) clean



$(BIN)/golint: | $(GOBASE)   # # # ❶
	@echo "  >  installing golint ..."
	#@-mkdir -p $(GOPATH)/src/golang.org/x/lint/golint
	#@cd $(GOPATH)/src/golang.org/x/lint/golint
	#@pwd
	#@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	#go get -v golang.org/x/lint/golint
	@echo "  >  installing golint ..."
	@$(GO) install golang.org/x/lint/golint
	@cd $(GOBASE)

$(BIN)/gocyclo: | $(GOBASE)  # # # ❶
	@echo "  >  installing gocyclo ..."
	@$(GO) install github.com/fzipp/gocyclo

$(BIN)/yolo: | $(GOBASE)     # # # ❶
	@echo "  >  installing yolo ..."
	@$(GO) install github.com/azer/yolo

$(BIN)/godoc: | $(GOBASE)     # # # ❶
	@echo "  >  installing godoc ..."
	@$(GO) install golang.org/x/tools/cmd/godoc

$(BASE):
	# @mkdir -p $(dir $@)
	# @ln -sf $(CURDIR) $@


## godoc: run godoc server at "localhost;6060"
godoc: | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	@echo "  $  cd $(GOPATH_) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	( cd $(GOPATH_) && pwd && godoc -v -index -http=:6060 -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps -goroot .; )
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b


## godoc1: run godoc server at "localhost;6060"
godoc1: # | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	#@echo "  $  GOPATH=$(GOPATH) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	godoc -v -index -http=:6060 -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps # -goroot $(GOPATH) 
	# gopkg.in/hedzr/errors.v2.New
	# -goroot $(GOPATH) -index
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b

## fmt: =`format`, run gofmt tool
fmt: format

## format: run gofmt tool
format: | $(GOBASE)
	@echo "  >  gofmt ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	gofmt -l -w -s .

## lint: run golint tool
lint: | $(GOBASE) $(GOLINT)
	@echo "  >  golint ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(GOLINT) ./...

## cov: =`coverage`, run go coverage test
cov: coverage

## gocov: =`coverage`, run go coverage test
gocov: coverage

## coverage: run go coverage test
coverage: | $(GOBASE)
	@echo "  >  gocov ..."
	@$(GO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=20m -test.short | tee coverage.log
	@$(GO) tool cover -html=coverage.txt -o cover.html
	@open cover.html

## coverage-full: run go coverage test (with the long tests)
coverage-full: | $(GOBASE)
	@echo "  >  gocov ..."
	@$(GO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=20m | tee coverage.log
	@$(GO) tool cover -html=coverage.txt -o cover.html
	@open cover.html

## codecov: run go test for codecov; (codecov.io)
codecov: | $(GOBASE)
	@echo "  >  codecov ..."
	@$(GO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic
	@bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN)

## cyclo: run gocyclo tool
cyclo: | $(GOBASE) $(GOCYCLO)
	@echo "  >  gocyclo ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(GOCYCLO) -top 20 .

## bench-std: benchmark test
bench-std:
	@echo "  >  benchmark testing ..."
	@$(GO) test -bench="." -run=^$ -benchtime=10s $(COVER_TEST_TARGETS)
	# go test -bench "." -run=none -test.benchtime 10s
	# todo: go install golang.org/x/perf/cmd/benchstat


## bench: benchmark test
bench:
	@echo "  >  benchmark testing (manually) ..."
	@$(GO) test ./fast -v -race -run 'TestQueuePutGetLong' -timeout=20m


## linux-test: call ci/linux_test/Makefile
linux-test:
	@echo "  >  linux-test ..."
	@-touch $(STDERR)
	@-rm $(STDERR)
	@echo $(MAKE) -f ./ci/linux_test/Makefile test 2> $(STDERR)
	@$(MAKE) -f ./ci/linux_test/Makefile test 2> $(STDERR)
	@echo "  >  linux-test ..."
	$(MAKE) -f ./ci/linux_test/Makefile all  2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/' 1>&2


## docker: docker build
docker:
	@if [ -n "$(DOCKER_APP_NAME)" ]; then \
	  echo "  >  docker build $(DOCKER_APP_NAME):$(VERSION)..."; \
	  docker build --build-arg CN=1 --build-arg GOPROXY="https://gocenter.io,direct" -t $(DOCKER_APP_NAME):latest -t $(DOCKER_APP_NAME):$(VERSION) . ; \
	else \
	  echo "  >  docker build not available since DOCKER_APP_NAME is empty"; \
	fi


## rshz: rsync to my TP470P
rshz:
	@echo "  >  sync to hz-pc ..."
	rsync -arztopg --delete $(GOBASE) hz-pc:$(HZ_PC_GOBASE)/src/github.com/hedzr/





.PHONY: directories

MKDIR_P = mkdir -p

directories: $(GOBIN)

$(GOBIN):
	$(MKDIR_P) $(GOBIN)




.PHONY: printvars info help all
printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

info:
	@echo "     GO_VERSION: $(GOVERSION)"
	@echo "        GOPROXY: $(GOPROXY)"
	@echo "         GOROOT: $(GOROOT) | GOPATH: $(GOPATH)"
	#@echo "    GO111MODULE: $(GO111MODULE)"
	@echo
	@echo "         GOBASE: $(GOBASE)"
	@echo "          GOBIN: $(GOBIN)"
	@echo "    PROJECTNAME: $(PROJECTNAME)"
	@echo "        APPNAME: $(APPNAME)"
	@echo "        VERSION: $(VERSION)"
	@echo "      BUILDTIME: $(BUILDTIME)"
	@echo "    GIT_VERSION: $(GIT_VERSION)"
	@echo "   GIT_REVISION: $(GIT_REVISION)"
	@echo
	@echo " MAIN_BUILD_PKG: $(MAIN_BUILD_PKG)"
	@echo "      MAIN_APPS: $(MAIN_APPS)"
	@echo
	#@echo "export GO111MODULE=on"
	@echo "export GOPROXY=$(GOPROXY)"
	#@echo "export GOPATH=$(GOPATH)"

all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

