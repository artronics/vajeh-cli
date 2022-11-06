.SILENT:

include .env/env

build:
	go build -o build/vajeh

vajeh: build
	./build/vajeh $(opt)

config-init: build
	./build/vajeh config init $(opt)

deploy: build
	./build/vajeh --config $(shell pwd)/.env/vajeh2.yaml deploy $(opt)

version: build
	./build/vajeh --config $(shell pwd)/.env/vajeh2.yaml version $(opt)

docker-login:
	docker login -u artronics -p $(DOCKER_HUB_TOKEN)

docker-build:
	docker build --platform amd64 -t artronics/vajeh-cli:local .

docker-push: docker-login
	docker push artronics/vajeh-cli:local

clean:
	rm -rf build vajeh-cli terraform.* .terraform terraform-test/.terraform* terraform-test/terraform.*

.PHONY: build run clean
