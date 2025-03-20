test:
	go test ./...

lint:
	golangci-lint run --fix

build-exporter:
	docker build -f Dockerfile_exporter -t artarts36/certmetrics-exporter:local .

run-exporter:
	docker run -v ./exporter:/app -p 8010:8010 artarts36/certmetrics-exporter:local
