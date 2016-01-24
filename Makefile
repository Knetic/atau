default: build
all: package

export GOPATH=$(CURDIR)/
export GOBIN=$(CURDIR)/.temp/

init: clean
	go get ./...

build: init
	go build -o ./.output/atau .

test:
	go test
	go test -bench=.

clean:
	@rm -rf ./.output/

fmt:
	@go fmt .
	@go fmt ./src/atau

dist: build test

	export GOOS=linux; \
	export GOARCH=amd64; \
	go build -o ./.output/atau64 .

	export GOOS=linux; \
	export GOARCH=386; \
	go build -o ./.output/atau32 .

	export GOOS=darwin; \
	export GOARCH=amd64; \
	go build -o ./.output/atau_osx .

	export GOOS=windows; \
	export GOARCH=amd64; \
	go build -o ./.output/atau.exe .

package: versionTest fpmTest dist

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(ATAU_VERSION) \
		-n atau \
		./.output/atau64=/usr/local/bin/atau \
		./docs/atau.7=/usr/share/man/man7/atau.7 \
		./autocomplete/atau=/etc/bash_completion.d/atau

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(ATAU_VERSION) \
		-n atau \
		-a i686 \
		./.output/atau32=/usr/local/bin/atau \
		./docs/atau.7=/usr/share/man/man7/atau.7 \
		./autocomplete/atau=/etc/bash_completion.d/atau

	@mv ./*.deb ./.output/

	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(ATAU_VERSION) \
		-n atau \
		./.output/atau64=/usr/local/bin/atau \
		./docs/atau.7=/usr/share/man/man7/atau.7 \
		./autocomplete/atau=/etc/bash_completion.d/atau
	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(ATAU_VERSION) \
		-n atau \
		-a i686 \
		./.output/atau32=/usr/local/bin/atau \
		./docs/atau.7=/usr/share/man/man7/atau.7 \
		./autocomplete/atau=/etc/bash_completion.d/atau

	@mv ./*.rpm ./.output/

fpmTest:
ifeq ($(shell which fpm), )
	@echo "FPM is not installed, no packages will be made."
	@echo "https://github.com/jordansissel/fpm"
	@exit 1
endif

versionTest:
ifeq ($(ATAU_VERSION), )

	@echo "No 'ATAU_VERSION' was specified."
	@echo "Export a 'ATAU_VERSION' environment variable to perform a package"
	@exit 1
endif
