
# dev: clean-kind ## Run local dev with Skaffold, watching for code changes. Deletes and recreates the test cluster.
# > kind create cluster --config=k8s/dev/skaffold/kind.yaml --name=ara-local-dev
# > skaffold debug -f skaffold.debug.yaml -p $(SKAFFOLD_PROFILE)
# .PHONY: dev

# Default Skaffold profile
SKAFFOLD_PROFILE ?= default
# HOME = $(HOME)

clean-kind: ## Deletes the local dev cluster created by Kind.
	kind delete cluster --name=lease-cluster
.PHONY: clean-kind

dev: clean-kind ## Run local dev with Skaffold, watching for code changes. Deletes and recreates the test cluster.
	kind create cluster --name=lease-cluster --config=cluster-config.yaml
	skaffold dev -p $(SKAFFOLD_PROFILE)
.PHONY: dev
