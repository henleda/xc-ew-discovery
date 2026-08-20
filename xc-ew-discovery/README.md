# xc-ew-discovery

Out-of-band east/west API discovery for Kubernetes, third-party gateways, and IBM DataPower. Prototype.

Start here:

1. `docs/ENVIRONMENT.md` before the first build. The dev box is arm64 and three things break because of it.
2. `make check-env` to confirm kernel and BTF.
3. `docs/MILESTONES.md` for what to build. M0 is a two-day spike and it comes before everything.

Working with Claude Code: `CLAUDE.md` loads every session and carries the hard constraints. `/milestone M1` starts a milestone with its definition of done in context. `/freeze-check` verifies nothing has drifted from the frozen schema.
