
.PHONY: dev
dev:
	docker container run \
		-ti \
		--rm \
		--mount type=bind,src=$$PWD,dst=/go/src/github.com/juli3nk/openfaas-connector-stan \
		--workdir /go/src/github.com/juli3nk/openfaas-connector-stan \
		--name stan-connector_dev \
		juli3nk/dev:go

.PHONY: build
build:
	docker image build \
		-t juli3nk/openfaas-connector-stan \
		.

.PHONY: push
push:
	docker image push juli3nk/openfaas-connector-stan
