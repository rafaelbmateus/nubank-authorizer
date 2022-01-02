# Nubank Authorizer

The `nubank-authorizer` is a micro-service designed to register account and transactions.
The service is consumed via the api interface, using json to traffic the payload objects.

| Verb |       URI       |      Action       |
|------|-----------------|-------------------|
| POST | `/accounts`     | Create an account |
| POST | `/transactions` | New transaction   |

## Project design

The project design follows the clean code principle with SOLID concepts.
The directory layout follows the
[golang-standards/project-layout](https://github.com/golang-standards/project-layout)
with some adaptations to make the project simpler, like [go](https://go.dev/).

## Get Started

This project considers that you have installed these tools:
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://www.docker.com)
- [Docker Compose](https://docs.docker.com/compose)

First, you can run all the tests.
Will ensure that everything is working.

```sh
$ make test

?       nubank/authorizer                 [no test files]
?       nubank/authorizer/database        [no test files]
ok      nubank/authorizer/database/memory 0.002s  coverage: 100.0% of statements
ok      nubank/authorizer/entity          0.001s  coverage: 100.0% of statements
ok      nubank/authorizer/handler         0.005s  coverage: 100.0% of statements
ok      nubank/authorizer/usecase         0.002s  coverage: 100.0% of statements
```

To run the service locally using docker and compose,
execute:

```sh
$ make up
```

Now, send a request to [localhost:3000](http://localhost:3000) using *get* verb.

```sh
$ curl localhost:3000
```

The result is:

```json
{
  "content": "nubank authorizer api"
}
```

To following container logs, execute:

```sh
$ make logs
```

## Create an account

To create a newaccount send a *post* to `/accounts`, for example:

```sh
$ curl -X POST http://localhost:3000/accounts \
  -d '{"account": {"activeCard": true, "availableLimit": 100}}'
```

## Register new transaction

To register a new transactions send a *post* to `/transactions`, like:

```sh
$ curl -X POST http://localhost:3000/transactions \
  -d '{"transaction": {"merchant": "ifood", "amount": 25, "time": "2020-12-22T10:00:00.000Z"}}'
```

## Basic Use Case

```sh
# create an account
$ curl -X POST http://localhost:3000/accounts \
  -d '{"account": {"activeCard": true, "availableLimit": 100}}'

{"account":{"activeCard":true,"availableLimit":100},"violations":[]}

# account-already-initialized rule
$ curl -X POST http://localhost:3000/accounts \
  -d '{"account": {"activeCard": true, "availableLimit": 350}}'

{"account":{"activeCard":true,"availableLimit":100},"violations":["account-already-initialized"]}

# create an transaction
$ curl -X POST http://localhost:3000/transactions \
  -d '{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}'

{"account":{"activeCard":true,"availableLimit":80},"violations":[]}

# doubled transaction rule
$ curl -X POST http://localhost:3000/transactions \
  -d '{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:02:00.000Z"}}'

{"account":{"activeCard":true,"availableLimit":80},"violations":["doubled-transaction"]}

# insuficient limit rule
$ curl -X POST http://localhost:3000/transactions \
  -d '{"transaction": {"merchant": "Burger King", "amount": 200, "time": "2019-02-13T10:00:00.000Z"}}'

{"account":{"activeCard":true,"availableLimit":80},"violations":["insufficient-limit"]}
```

## Helper

execute `make help` to show details about makefile tasks:

```sh
$ make help
usage: make [target]

development:
  build                           Build docker images.
  up                              Run containers in detach.
  restart                         Restart development environment.
  stop                            Stop development environment and remove containers orphans.
  logs                            Follows development logs.
  shell                           Start a shell session within the container.

lint:
  lint                            Run static analysis code.

other:
  one-shot                        Execiute basic use case using curl.
  help                            Show this help.

test:
  test                            Run the tests.
  coverage                        Generate coverage files.
```

## Challenge description

For more information about the challenge, [see here](docs/challenge.md)
