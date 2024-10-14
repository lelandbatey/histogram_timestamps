
all: histogram_timestamps

histogram_timestamps: bundle.js index.html main.go tbin/tbin.go timeformat/timeformat.go
	go build

bundle.js: ./jsbuild/node_modules
	./jsbuild/build_bundle_for_makefile.sh

./jsbuild/node_modules:
	cd ./jsbuild/; npm install

.PHONY: clean
clean:
	rm -f ./histogram_timestamps ./bundle.js ./jsbuild/bundle.js templates.go

.PHONY: install
install: all
	go install

