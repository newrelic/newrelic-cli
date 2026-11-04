# Security Policy

## Supported Versions

We release security patches for the latest major version of the New Relic CLI.

| Version | Supported |
|---------|-----------|
| Latest  | ✅        |
| Older   | ❌        |

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Report vulnerabilities via New Relic's coordinated disclosure program:
- Email: security@newrelic.com
- Bug bounty: https://hackerone.com/newrelic

We aim to acknowledge reports within 2 business days and to provide a remediation timeline within 7 business days.

---

## Release Pipeline Security Controls

The `newrelic-cli` release pipeline meets the requirements of New Relic's
supply-chain security policy ("Securing what we publish to customers", Nov 2026).

### Credential management (Req 3)

AWS credentials are short-lived OIDC tokens issued per run. No long-lived
`AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` secrets are present in the
release path. The IAM role (`arn:aws:iam::017663287629:role/virtuoso`) trust
policy requires:

```
repo:newrelic/newrelic-cli:ref:refs/tags/v*
```

Tokens are scoped to this exact repository and tag pattern — no other workflow,
branch, or fork can obtain the role.

### Two-step release process (Req 3)

Releases follow a mandatory A1 → A2 split:

| Workflow | Trigger | Purpose |
|---|---|---|
| `release-prepare.yml` (A1) | Manual `workflow_dispatch` from `main`, gated by the **production** GitHub Environment | Version bump, CHANGELOG update, protected tag creation |
| `release-publish.yml` (A2) | Automatic on `push: tags: v*` | Build, sign, upload to S3 / GitHub / Docker / Snap |

### Environment gate (Req 2)

A1 is gated by the **production** GitHub Environment, which requires:
- At least one reviewer from the `nr-developer-toolkit` team
- Self-review is blocked (the person who triggers cannot be their own approver)
- Deploys from `main` only

### Protected tag ruleset (Req 2)

The `v*` tag namespace is protected by a GitHub tag ruleset. Only the service
account used by A1 (`RELEASE_TOKEN`) can push `v*` tags. Direct tag pushes by
individual contributors are blocked. Every bypass is audited by GitHub.

### Action pinning (Req 4)

All GitHub Actions in release workflows are pinned to immutable commit SHAs,
not floating version tags. Go toolchain binaries are installed at specific
semver versions (not `@latest`). See individual workflow files for SHA comments.

### Build provenance (Req 5)

Docker images are built with provenance attestation enabled (BuildKit default).
Binary checksums (`checksums_sha256.txt`) and a detached PGP signature
(`checksums_sha256.txt.sig`) are published to every GitHub release. The signing
key fingerprint is `4F9A9B5B96EC30B9`.

### Installer integrity (Req 7)

`scripts/install.sh` downloads `checksums_sha256.txt` from the same release
path as the binary tarball and verifies the SHA-256 checksum before unpacking
or executing anything.

---

## Emergency Release Path (Break-Glass)

If the A1/A2 pipeline is unavailable (e.g., GitHub Actions outage, OIDC
federation failure), a release can be performed using the legacy
`release.yml` workflow. This workflow uses long-lived credentials and should
be used **only as a last resort**. It is documented as deprecated in its file
header.

Steps:
1. Confirm the pipeline outage is real and not a transient error.
2. Get approval from the team lead or on-call engineer.
3. Manually trigger `release.yml` from the GitHub Actions UI.
4. File an incident noting the use of the break-glass path and the reason.
5. Rotate `AWS_ACCESS_KEY_ID_V2` / `AWS_SECRET_ACCESS_KEY_V2` after use
   if there is any suspicion of exposure.

## Credential Revocation

| Credential | Revocation path |
|---|---|
| AWS OIDC role | IAM Console → Roles → `virtuoso` → Trust policy: remove the `newrelic-cli` repo condition |
| `AWS_ACCESS_KEY_ID_V2` (break-glass only) | IAM Console → Users → `virtuoso_user` → Security credentials → Deactivate/delete |
| `RELEASE_TOKEN` (GitHub PAT) | GitHub → Settings → Developer settings → Personal access tokens → Revoke |
| `PGP_PRIVATE_KEY` | Revoke via the public keyserver: `gpg --keyserver keys.openpgp.org --send-keys 4F9A9B5B96EC30B9` after uploading a revocation certificate |
| `DOCKER_PASSWORD` | Docker Hub → Account settings → Security → Access tokens |

After revoking any credential, update the corresponding GitHub Actions secret
and re-test the release pipeline with a dry-run before the next real release.

## Rollback Process

If a release is found to be compromised or broken after publishing:

1. **GitHub release**: Edit the release on GitHub, mark it as pre-release or
   delete it. This removes the release assets from the GitHub release page but
   does **not** remove objects from S3.

2. **S3 artifacts**: Use the AWS console or CLI to delete the version-specific
   prefix (`s3://nr-downloads-main/install/newrelic-cli/<bad-version>/`).

3. **`currentVersion.txt`**: Overwrite with the last known-good version:
   ```bash
   echo "v<last-good-version>" | aws s3 cp - \
     s3://nr-downloads-main/install/newrelic-cli/currentVersion.txt \
     --cache-control no-cache
   ```
   This prevents new `install.sh` runs from pulling the bad version.

4. **Docker images**: Delete the affected tags from Docker Hub via the UI or
   `docker hub delete-tag` CLI.

5. **Notify users**: Post a security advisory on the GitHub repository.
