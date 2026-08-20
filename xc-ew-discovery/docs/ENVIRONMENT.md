# Environment

Development target: one amd64 cloud VM running Ubuntu 24.04 with k3s and DataPower for Developers side by side. A managed EKS cluster gets created on demand for M6 and for the node-exclusion evidence, then torn down.

## Why one box carries M0 through M5

Everything through M5 is answerable on a single node. Cross-node traffic and the coverage math start mattering at M4, and adding a second and third k3s agent in the same VPC covers that for a few dollars a day.

amd64 removes two problems the arm64 path carried. Beyla's Go uprobe support is best tested on amd64, and DataPower ships no arm64 build at all.

## Sizing

8 vCPU and 32 GB is comfortable. 4 vCPU and 16 GB works with DataPower running lean.

Budget roughly: DataPower for Developers wants 4 GB and 2 cores. The k3s control plane, demo app, and Beyla take another 4 GB. M0 adds two JVMs on top.

- AWS `m7i.2xlarge` with an Ubuntu 24.04 AMI. Stop the instance when idle.
- GCP `n2-standard-8` with an Ubuntu 24.04 image.
- Hetzner CCX33. Cheapest by a wide margin, and the kernel is fully yours.

## Do not develop this on a Mac

Docker Desktop and kind on macOS run Linux inside a VM. uprobe behavior in that VM does not represent a customer node, so an M0 pass there proves nothing.

## Kernel and BTF

Ubuntu 24.04 ships 6.8 with BTF compiled in. Confirm anyway before M1.

```
make check-env
```

Avoid images with custom kernels lacking `CONFIG_DEBUG_INFO_BTF`. Ubuntu 22.04 and newer and Amazon Linux 2023 are safe. Container-Optimized OS on GKE has a read-only root filesystem that complicates eBPF tooling, so pick Ubuntu node images there.

## DataPower

IBM ships a free developer edition as a container, DataPower Gateway for Developers, on the IBM Container Registry. It is throughput-limited and it exposes the SOMA and REST management interfaces, which is everything M5 needs.

Run it on the dev VM alongside k3s rather than inside the cluster. The connector is a network client, and the appliance is not a Kubernetes workload in a real deployment either.

Ports: 5550 SOMA, 5554 REST management, 9090 web GUI. First boot needs an interactive start to accept the license and set an admin password.

Commit sanitized responses into `scrape/datapower/testdata/` as you go, so tests run without the container.

## Managed cluster validation

Create on demand, tear down after. Two things need a cluster you do not control and neither needs it for long.

1. M6, confirming the CE deploys and the egress assertion holds outside a hand-built environment.
2. Node exclusion evidence for the spec. EKS Fargate refuses privileged DaemonSets and GKE Autopilot constrains them. Capture the actual error messages. They belong in the install docs and in field enablement, and they are the honest answer to "what percentage of the fleet do we lose."

Use EKS with a managed node group on Ubuntu or AL2023 AMIs. A Graviton node group is worth one pass eventually, since arm64 nodes are common in enterprise EKS now, but treat it as a GA concern rather than a prototype one.

## k3s

```
make dev-cluster
make deploy-demo
make deploy-sensor
```

k3s uses containerd, not Docker. Beyla needs the containerd cgroup path and `demo/beyla-values.yaml` sets it. If observations arrive with an empty `Workload.Pod`, check that first.

Adding nodes for the M4 coverage math:

```
curl -sfL https://get.k3s.io | K3S_URL=https://<server>:6443 K3S_TOKEN=<token> sh -
```
