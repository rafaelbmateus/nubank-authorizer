.DEFAULT_GOAL = help
compose = docker-compose -f build/docker-compose.yml

.PHONY: build
build: ##@development Build docker images.
	$(compose) build

.PHONY: up
up: ##@development Run containers in detach.
	$(compose) up -d

.PHONY: restart
restart: ##@development Restart development environment.
	$(compose) restart

.PHONY: stop
stop: ##@development Stop development environment and remove containers orphans.
	$(compose) down -v --remove-orphans

.PHONY: logs
logs: ##@development Follows development logs.
	$(compose) logs -f --tail=100

.PHONY: shell
shell: ##@development Start a shell session within the container.
	$(compose) run --rm app /bin/sh

lint_version ?= v1.40-alpine
.PHONY: lint
lint: ##@lint Run static analysis code.
	docker run --rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:$(lint_version) \
		golangci-lint run --tests=false --timeout 3m

.PHONY: test
test: ##@test Run the tests.
	$(compose) run --rm app go test -coverprofile coverage.out ./...

.PHONY: coverage
coverage: ##@test Generate coverage files.
	$(compose) run --rm app go tool cover -html coverage.out -o coverage.html

.PHONY: one-shot
one-shot: ##@other Execiute basic use case using curl.
	@curl http://localhost:3000 && echo
	@curl -X POST http://localhost:3000/accounts -d '{"account": {"activeCard": true, "availableLimit": 100}}' && echo
	@curl -X POST http://localhost:3000/accounts -d '{"account": {"activeCard": true, "availableLimit": 350}}' && echo
	@curl -X POST http://localhost:3000/transactions -d '{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}' && echo
	@curl -X POST http://localhost:3000/transactions -d '{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:02:00.000Z"}}' && echo
	@curl -X POST http://localhost:3000/transactions -d '{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:02:01.000Z"}}' && echo
	@curl -X POST http://localhost:3000/transactions -d '{"transaction": {"merchant": "Burger King", "amount": 90, "time": "2019-02-13T10:00:00.000Z"}}'

GREEN  = $(shell tput -Txterm setaf 2)
WHITE  = $(shell tput -Txterm setaf 7)
YELLOW = $(shell tput -Txterm setaf 3)
RESET  = $(shell tput -Txterm sgr0)
HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
	print "usage: make [target]\n\n"; \
	for (sort keys %help) { \
	print "${WHITE}$$_:${RESET}\n"; \
	for (@{$$help{$$_}}) { \
	$$sep = " " x (32 - length $$_->[0]); \
	print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
	}; \
	print "\n"; }

.PHONY: help
help: ##@other Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

