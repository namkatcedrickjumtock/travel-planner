install-tools:
	if [ ! $$(which go) ]; then \
		echo "goLang not found."; \
		echo "Try installing go..."; \
		exit 1; \
	fi
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/axw/gocov/gocov@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1
	go get golang.org/x/tools/cmd/goimports
	go install github.com/AlekSi/gocov-xml@latest
	if [ ! $$( which migrate ) ]; then \
		echo "The 'migrate' command was not found in your path. You most likely need to add \$$HOME/go/bin to your PATH."; \
		exit 1; \
	fi
	

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

create-migration: ## usage: make name=new create-migration
	migrate create -ext sql -dir ./db/migrations -seq $(name)
 

run:database
	go mod tidy
	if [ ! -f '.env' ]; then \
		cp .env.example .env; \
	fi
	go run ./cmd/api/api.go 

database:
	docker-compose up -d
	