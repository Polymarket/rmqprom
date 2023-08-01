#######################################################
############## formats, lint, and tests ###############
#######################################################

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	gofmt -s -w ./.

.PHONY: lint
lint: 
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Linting code..."
	@echo "----------------------------------------------------------------"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3 run ./... -E gofmt --config=.golangci.yaml 
	@echo "Linting complete!"

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Testign the code..."
	@echo "----------------------------------------------------------------"
	GOPRIVATE=${PRIVATE_REPOS} go test ./... -v 
	@echo "Tests complete!"
