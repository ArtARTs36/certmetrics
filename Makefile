test:
	go test ./...

lint:
	golangci-lint run --fix

build-exporter:
	docker build -f Dockerfile -t artarts36/certmetrics:local .

run-exporter:
	docker run -v ./exporter:/app -p 8010:8010 artarts36/certmetrics:local

check: test lint
