
.PHONY: build
build:
	docker image build -t juli3nk/openfaas-connector-stan .

.PHONY: push
push:
	docker image push juli3nk/openfaas-connector-stan
