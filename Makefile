DIST_DIRS := find * -type d -exec
VERSION := $(shell git describe --tags)

.PHONY: build
build:
	go build -o sqspush -ldflags "-X main.version=${VERSION}" sqspush.go

.PHONY: install
install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./sqspush ${DESTDIR}/usr/local/bin/sqspush

.PHONY: integration-test
integration-test:
	go build
	./sqspush

.PHONY: clean
clean:
	rm -f ./sqspush
	rm -rf ./dist

.PHONY: bootstrap-dist
bootstrap-dist:
	go get -u github.com/franciscocpg/gox
	cd ${GOPATH}/src/github.com/franciscocpg/gox && git checkout dc50315fc7992f4fa34a4ee4bb3d60052eeb038e
	cd ${GOPATH}/src/github.com/franciscocpg/gox && go install

.PHONY: build-all
build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows freebsd openbsd netbsd" \
	-arch="amd64 386 armv5 armv6 armv7 arm64" \
	-osarch="!darwin/arm64" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

.PHONY: dist
dist: build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf sqspush-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r sqspush-${VERSION}-{}.zip {} \; && \
	cd ..
