
all: histogram_timestamps

histogram_timestamps: templates.go
	go build

templates.go: bundle.js index.html external_deps
	go-bindata -o templates.go index.html bundle.js

bundle.js: ./jsbuild/node_modules
	./jsbuild/build_bundle_for_makefile.sh

./jsbuild/node_modules:
	cd ./jsbuild/; npm install

.PHONY: external_deps
external_deps:
	go get -u github.com/kevinburke/go-bindata/go-bindata

.PHONY: clean
clean:
	rm -f ./histogram_timestamps ./bundle.js ./jsbuild/bundle.js templates.go

.PHONY: install
install: all
	go install

