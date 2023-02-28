CURDIR=$(shell pwd)
APIDIR=$(CURDIR)/api

build-all:
	cd checkout && GOOS=linux make build
	cd loms && GOOS=linux make build
	cd notifications && GOOS=linux make build

run-all: build-all
	docker compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

generate-all:
	cd checkout && make generate APIDIR=$(APIDIR)
	cd loms && make generate APIDIR=$(APIDIR)

