# Beta Release Workflow

**Trigger:** Push tagu `v*-beta.*` (e.g. `v1.1.0-beta.1`)

## What It Does

1. **Build binaries** — same matrix as stable (linux/darwin/windows, amd64/arm64), named with beta version
2. **GitHub prerelease** — creates GitHub release with `prerelease: true`
3. **Docker image** — pushes `tedyno/ticktock-mcp:<version>` only, **no `latest` tag**

## Key Differences from Stable

| | Stable | Beta |
|---|---|---|
| Trigger | `v*` (excluding beta) | `v*-beta.*` |
| Docker `latest` | yes | **no** |
| GitHub prerelease | no | **yes** |

## Changes Required

1. Create `.github/workflows/beta-release.yml` — new workflow for beta releases
2. Modify `.github/workflows/release.yml` — exclude beta tags from trigger so both workflows don't fire on the same tag
