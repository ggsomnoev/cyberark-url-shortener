.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	ginkgo ./...	

.PHONY: run-docker
run-docker:
	docker compose down -v --remove-orphans 
	docker compose up --build 