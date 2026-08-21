GOARCH ?= amd64

.PHONY: build test fmt vet dev-cluster deploy-demo deploy-sensor uninstall check-env clean

build:
	GOARCH=$(GOARCH) go build ./...

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

check-env:
	@echo "kernel: $$(uname -r)"
	@echo "arch:   $$(uname -m)"
	@test -f /sys/kernel/btf/vmlinux && echo "btf:    present" || echo "btf:    MISSING, see docs/ENVIRONMENT.md"

dev-cluster:
	curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable traefik" sh -
	@echo "kubeconfig at /etc/rancher/k3s/k3s.yaml"

deploy-demo:
	kubectl apply -f demo/apps/

deploy-sensor:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm upgrade --install beyla grafana/beyla -n beyla --create-namespace -f demo/beyla-values.yaml

# uninstall removes the sensor and asserts a clean node. P0#12, M9.
uninstall:
	helm uninstall beyla -n beyla || true
	@echo "P0#12 clean-node check:"
	@ls /sys/fs/bpf 2>/dev/null | grep -i beyla && echo "FAIL: pinned BPF objects remain" || echo "clean: no pinned BPF objects"

clean:
	rm -rf ./specs ./bin
