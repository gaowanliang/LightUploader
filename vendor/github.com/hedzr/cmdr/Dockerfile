ARG BASE_BUILD_IMAGE
FROM ${BASE_BUILD_IMAGE:-golang:1.15} AS builder

ARG VERSION
ARG W_PKG
ARG GO111MODULE
ARG GOPROXY
ARG CN
ARG WORKDIR
ARG APPNAME

# docker build --build-arg CN=1 -t awesome-tool:latest . 

ENV AN=${APPNAME:-fluent}
ENV SRCS=./examples/fluent
ENV WDIR=${WORKDIR:-/var/lib/$AN}
ENV GIT_REVISION	""
ENV GOVERSION			"1.15"
ENV BUILDTIME			""
ENV LDFLAGS				""

WORKDIR /go/src/github.com/hedzr/$AN/
COPY    .    .
RUN ls -ls ./; \
		W_PKG=${W_PKG:-github.com/hedzr/cmdr/conf}; \
		GOPROXY=${GOPROXY:-https://goproxy.io,direct}; \
		V1=$(grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+"); \
		VERSION=${VERSION:-$V1}; \
		GIT_REVISION="$(git rev-parse --short HEAD)"; \
		GOVERSION="$(go version)"; \
		BUILDTIME="$(date -u '+%Y-%m-%d_%H-%M-%S')"; \
		LDFLAGS="-s -w \
			-X '$W_PKG.Githash=$GIT_REVISION' \
			-X '$W_PKG.GoVersion=$GOVERSION' \
			-X '$W_PKG.Buildstamp=$BUILDTIME' \
			-X '$W_PKG.Version=$VERSION'"; \
		echo;echo;echo "Using GOPROXY: $GOPROXY";echo "    CN: $CN";echo; \
		CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo \
			-ldflags "$LDFLAGS" \
			-o bin/$AN $SRCS && \
		ls -la bin/






ARG BASE_IMAGE
FROM ${BASE_IMAGE:-alpine:latest}

ARG CN
ARG VERSION
ARG WORKDIR
ARG CONFDIR
ARG APPNAME

ENV AN=${APPNAME:-fluent}
ENV SRCS=./examples/fluent
ENV WDIR=${WORKDIR:-/var/lib/$AN}
ENV CDIR=${CONFDIR:-/etc/$AN}

LABEL by="hedzr" \
			version="$VERSION" \
			com.hedzr.cmdr-fluent.version="$VERSION" \
			com.hedzr.cmdr-fluent.release-date="$(date -u '+%Y-%m-%d_%H-%M-%S')" \
			description="awesome-tool a command-line tool to retrieve the stars of all repos in an awesome-list"

COPY --from=builder /go/src/github.com/hedzr/$AN/ci/etc/$AN /etc/$AN

RUN ls -la $CDIR/ $CDIR/conf.d && echo "    CN: $CN"; \
    [[ "$CN" != "" ]] && { \ 
      cp /etc/apk/repositories /etc/apk/repositories.bak; \
      echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories; \
      echo;echo;echo "apk updating...";apk update; }; \
    apk --no-cache add ca-certificates && \
    mkdir -p $WDIR/output /var/log/$AN /var/run/$AN && \
    ls -la $WDIR/output /var/log/$AN /var/run/$AN
    

VOLUME  [	"$WDIR/output", "$CDIR/conf.d" ]
WORKDIR $WDIR
COPY --from=builder /go/src/github.com/hedzr/$AN/bin/$AN .
RUN echo $WDIR && echo $AN && ls -la $WDIR

ENTRYPOINT [ "./fluent" ]
CMD [ "--help" ]
