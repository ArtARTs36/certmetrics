test:
	go test ./...

lint:
	golangci-lint run --fix

build-exporter:
	docker build -f Dockerfile -t artarts36/certmetrics:local .

run-exporter: build-exporter
	docker run -v ./:/app -p 8010:8010 artarts36/certmetrics:local

check: test lint
