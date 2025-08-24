include .env

LOCAL_BIN := $(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.15.1


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
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore
