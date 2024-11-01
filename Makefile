#!make

main_bin_name = go-dynamicdns
cmd_path = ./cmd/go-dynamicdns

.PHONY: clean
clean:
	if [ -d ./tmp ]; then rm -r ./tmp; fi


.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...


.PHONY: build/dev
build/dev: clean
	go build -o=./tmp/bin/${main_bin_name} ${cmd_path}

.PHONY: run/dev
run/dev: build/dev
	./tmp/bin/${main_bin_name}
