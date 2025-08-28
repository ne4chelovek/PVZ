include .env

LOCAL_BIN := $(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.15.1
	GOBIN=$(LOCAL_BIN) go install github.com/itchyny/gojq/cmd/gojq@latest
	GOBIN=$(LOCAL_BIN) go install github.com/tsenart/vegeta@latest

migrate-up:
	goose -dir migrations/ postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations/ postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations/ postgres "$(DB_URL)" status

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=PVZ/internal/service/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=PVZ/internal/service/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out
	go tool cover -func=./coverage.out | grep "total"
	@grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore


.PHONY: load-test clean

BASE_URL      ?= http://localhost:8080
DUMMY_LOGIN   := $(BASE_URL)/dummyLogin
LOCAL_BIN     := ./bin

load-test:
	$(eval TOKEN := $(shell curl -s -X POST $(DUMMY_LOGIN) -H "Content-Type: application/json" -d '{"role":"moderator"}' | $(LOCAL_BIN)/gojq -r '.token'))
	@echo " Load Test: 100 RPS, 10s, POST /pvz"
	@mkdir -p tmp
	printf 'POST $(BASE_URL)/pvz' | $(LOCAL_BIN)/vegeta attack \
		-rate=100 \
		-duration=10s \
		-timeout=10s \
		-header="Authorization: Bearer $(TOKEN)" \
		-header="Content-Type: application/json" \
		-body=testdata/request-body.json | \
	$(LOCAL_BIN)/vegeta encode > tmp/vegeta.bin

	$(LOCAL_BIN)/vegeta report tmp/vegeta.bin

negative-test:
	$(eval TOKEN := $(shell curl -s -X POST $(DUMMY_LOGIN) -H "Content-Type: application/json" -d '{"role":"moderator"}' | $(LOCAL_BIN)/gojq -r '.token'))
	@echo " Negative test: invalid city"
	@printf 'POST $(BASE_URL)/pvz' | $(LOCAL_BIN)/vegeta attack \
		-rate=50 \
		-duration=5s \
		-header="Authorization: Bearer $(TOKEN)" \
		-header="Content-Type: application/json" \
		-body=testdata/negative-body.json | \
	$(LOCAL_BIN)/vegeta report | grep -E "(Requests|Latencies|Success|Status)"

clean:
	rm -f /tmp/pvz_test_token.txt