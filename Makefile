clean-kind: ## Deletes the local dev cluster created by Kind.
	kind delete cluster --name=lease-cluster

run-skaffold: ## Run just skaffold
	skaffold dev -p default

dev: clean-kind ## Run local dev with Skaffold, watching for code changes. Deletes and recreates the test cluster.
	kind create cluster --name=lease-cluster --config=cluster-config.yaml
	$(MAKE) run-skaffold

.PHONY: dev clean-kind run-skaffold
