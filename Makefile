test:
	go test ./...
	cd ./pkg/collector && go test ./...
	cd ./pkg/jwtm && go test ./...
	cd ./pkg/x509m && go test ./...

lint:
	golangci-lint run --fix

build-exporter:
	docker build -f Dockerfile -t artarts36/certmetrics:local .

run-exporter: build-exporter
	docker run -v ./:/app -p 8010:8010 artarts36/certmetrics:local

check: test lint
