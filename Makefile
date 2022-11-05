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

plan: build
	./build/vajeh --config $(shell pwd)/.env/vajeh2.yaml deploy --dryrun $(opt)

config-init: build
	./build/vajeh config init $(opt)

docker-login:
	docker login -u artronics -p $(DOCKER_HUB_TOKEN)

docker-build:
	docker build --platform amd64 -t artronics/vajeh-cli:local .

docker-push: docker-login
	docker push artronics/vajeh-cli:local

clean:
	rm -rf build vajeh-cli terraform.* .terraform terraform-test/.terraform* terraform-test/terraform.*

.PHONY: build run clean
