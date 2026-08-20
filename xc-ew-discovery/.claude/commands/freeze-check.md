---
description: Verify the frozen schema and the hard constraints have not drifted
---

Check the following and report any violation with the file and line.

1. `model/event.go` matches the field set recorded in `docs/DECISIONS.md` D3. Any addition, removal, rename, or type change without a corresponding decision entry is a violation.
2. No `ClassifierHit` value-bearing field exists anywhere in the repo.
3. No code path stores, logs, or transmits a request or response body outside `infer/`.
4. `collector/egress.go` has no bypass, debug flag, or verbose mode that skips `AssertNoBodies`.
5. No `Workload` field is ever assigned an IP address.
6. No eBPF program source exists in this repo. The sensor is adopted, not written.
7. No top-level package beyond: model, sensor, enrich, infer, scrape, collector, cmd.

Report violations only. Say nothing about the parts that pass.
