.SILENT:

include .env/env

build:
	go build -o build/vajeh

run: build
	./build/vajeh $(opt)

dev: build
	./build/vajeh --config $(shell pwd)/.env/vajeh.yaml $(opt_dev)

deploy: build
	./build/vajeh --config $(shell pwd)/.env/vajeh2.yaml deploy $(opt)

config-init: build
	./build/vajeh config init $(opt)

.PHONY: build run
