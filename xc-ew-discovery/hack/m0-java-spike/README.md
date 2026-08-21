# M0. Java TLS spike

Two days. Do this before writing anything else in the repo.

The JVM terminates TLS in SunJSSE, not OpenSSL, so uprobes on libssl miss every Java-to-Java mTLS call. The target accounts for this product are the most Java-heavy accounts in existence. This is the highest-risk unknown in the program.

Throwaway code. Repo conventions do not apply here.

## Setup

Two Spring Boot services on the k3s cluster, one calling the other over mTLS with a self-signed CA. Run Beyla against them with the values in `demo/beyla-values.yaml`.

## What to record

Write the answers into `docs/DECISIONS.md` under "M0 result".

1. Does Beyla produce any observation at all for the Java-to-Java call, even without payload.
2. Does it produce L7 detail: method, path, status.
3. Does it produce payload shape.
4. Whether the answer changes on arm64. Out of scope for this spike: the dev box is amd64 (D8), so run the spike there. arm64 gets one validation pass on a Graviton EKS node group before GA, not here.
5. If payload capture fails, what a JVMTI agent costs in install friction. A second agent inside the JVM is a different approval conversation than a node DaemonSet, so price that honestly.

## Decision this drives

Java capture stays P1 and ships later, or it becomes a blocker for financial services accounts and changes the target segment for v1. Do not start M1 until this is written down.
