# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `kubevault.dev/cli` — the [KubeVault](https://kubevault.com/) command-line, distributed as a **`kubectl` plugin** (binary name `kubectl-vault`). Used to drive KubeVault out-of-band operations: VaultPolicy/VaultPolicyBinding approval, unseal-key/root-token recovery, GCP/AWS/Azure credential generation, and merging restic-backed snapshot secrets.

Top-level subcommands (from `pkg/cmds/root.go`):

- `approve`, `deny` — handle VaultPolicy / VaultPolicyBinding approval flow.
- `revoke` — revoke credentials issued by a VaultRole.
- `generate` — generate credentials via a configured secrets engine.
- `root-token` — read/rotate the Vault root token.
- `unseal-key` — read/manage Vault unseal keys (typically pulled from a `TokenKeysStore`).
- `merge-secrets` — merge restic-backed secret snapshots.

Plus `version` and `completion`.

## Architecture

- `cmd/kubectl-vault/main.go` — entry point; calls into `pkg/cmds`.
- `pkg/cmds/` — Cobra command tree. One file per subcommand.
- `pkg/generate/` — credential-generation logic shared by `generate` subcommand variants (AWS, Azure, GCP, database engines).
- `pkg/token-keys-store/` — abstraction for reading/storing Vault unseal keys and root tokens; backed by Kubernetes Secrets, cloud KMS, or external storage depending on the operator's configuration.
- `hack/`, `Makefile` — AppsCode build harness (everything runs inside `ghcr.io/appscode/golang-dev`). Binary builds for **5 platforms**: linux amd64/arm/arm64 plus `windows/amd64`, `darwin/amd64`, `darwin/arm64` (kubectl plugins need to run on operator workstations).
- `vendor/` — checked-in deps.

The binary uses the conventional `kubectl-vault` name so it auto-attaches under `kubectl vault` once on `$PATH`. There is no Docker image — this is a host CLI.

API types come from `kubevault.dev/apimachinery`.

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make ci` — CI pipeline.
- `make build` — build for host OS/ARCH into `bin/<os>_<arch>/kubectl-vault`.
- `make all-build` — build for every `BIN_PLATFORMS` (linux amd64/arm/arm64 + windows/amd64 + darwin/amd64 + darwin/arm64).
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make add-license` / `make check-license` — manage license headers.

There is **no container target** — this CLI does not ship as an image.

Run a single Go test (requires a local Go toolchain):

```
go test ./pkg/... -run TestName -v
```

To use the CLI locally after building:

```
PATH=$PWD/bin/<os>_<arch>:$PATH kubectl vault --help
```

## Conventions

- Module path is `kubevault.dev/cli` (vanity URL). Imports must use that.
- License: `LICENSE.md` (AppsCode); use `make add-license` to apply headers to new files.
- Sign off commits (`git commit -s`); contributions follow the DCO.
- Vendor directory is checked in — keep `go mod tidy && go mod vendor` clean.
- Binary name is `kubectl-vault` so kubectl picks it up as a plugin; do not rename without also moving the `cmd/` directory.
- New subcommand: add a `pkg/cmds/<name>.go`, register it in `root.go`'s `NewCmd()`. Cloud-provider credential generators belong in `pkg/generate/<provider>.go`, not as parallel `pkg/cmds/` files.
- The `token-keys-store` package abstracts unseal-key storage; new backends (additional cloud KMS, etc.) implement that interface and don't leak into the cmd surface.
- Builds linux/windows/darwin host binaries; do not pull in linux-only or cgo deps.
- API types come from `kubevault.dev/apimachinery`; do not redefine them here.
