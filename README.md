# ScopeSentry-Scan
The scanning component of ScopeSentry.

## Documentation and Installation

- Website: https://www.scope-sentry.top/
- Main project: https://github.com/Autumn-27/ScopeSentry

## Docker release workflow

On published GitHub Releases, `.github/workflows/docker-release.yml` builds `./dockerfile` with Buildx and pushes to Docker Hub as `${DOCKERHUB_USERNAME}/scopesentry-scan`.

Required repository secrets:
- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

Tag behavior:
- Always pushes the release tag (`github.event.release.tag_name`)
- Pushes `latest` only when the release is not a prerelease
