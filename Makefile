test:
	go test ./...

lint:
	golangci-lint run --fix

build-exporter:
	docker build -f Dockerfile_exporter -t artarts36/certmetrics-exporter:local .
