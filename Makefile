.PHONY: docker
docker:
	@rm mini-book || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o mini-book .
	@docker rmi cyanaqing/mini-book:v0.0.1
	@docker build -t cyanaqing/mini-book:v0.0.1 .