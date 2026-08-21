# xc-ew-discovery

Out-of-band east/west API discovery for Kubernetes, third-party gateways, and IBM DataPower. Prototype.

Start here:

1. `docs/ENVIRONMENT.md` before the first build. The dev box is an amd64 cloud VM running k3s; do not build on macOS or arm64 (see `docs/DECISIONS.md` D8).
2. `make check-env` to confirm kernel and BTF.
3. `docs/MILESTONES.md` for what to build. M0 is a two-day spike and it comes before everything.

CI: `.github/workflows/ci.yml` runs build, vet, gofmt, and tests, plus shell guards for the hard constraints (no eBPF source, no `ClassifierHit` value field, no new top-level package) on every push and pull request. The no-body-egress test (constraint 2) joins it when `collector/egress.go` is implemented at M6.

Working with Claude Code: `CLAUDE.md` loads every session and carries the hard constraints. `/milestone M1` starts a milestone with its definition of done in context. `/freeze-check` verifies nothing has drifted from the frozen schema.
