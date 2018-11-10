
BINARY_NAME=mars
MAIN_FILE=main.go
RUNTIME_MODE=dev

build:
	go build $(RACE) -o bin/${BINARY_NAME} ${MAIN_FILE}

build-race: enable-race build

run: build
	./bin/${BINARY_NAME} server --env ${RUNTIME_MODE}

run-race: enable-race run

test:
	go test $(RACE) ./...

test-race: enable-race test

enable-race:
	$(eval RACE = -race)

package: build-vue statik
	bash ./package.sh

package-all: build-vue statik
	bash ./package.sh -p 'linux darwin windows'

build-vue:
	cd web/vue && yarn run build
	cp -r web/vue/dist/* web/public/

install-vue:
	cd web/vue && yarn install

run-vue:
	cd web/vue && yarn run dev

statik:
	go get github.com/rakyll/statik
	go generate ./...

clean:
	rm bin/${BINARY_NAME}

.PHONY: clean statik run-vue install-vue build-vue package-all package enable-race
.PHONY: test-race test build build-race run run-race