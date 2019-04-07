## Description

Create Github nightly tags from CI pipeline.

## Output variables
```console
none
```

## Configuration defaults
```console
GOOPSC_GITHUB_USER=
GOOPSC_GITHUB_TOKEN=
GOOPSC_GIT_USER_EMAIL=travis@travis-ci.org
GOOPSC_GIT_USER_NAME=Travis CI
```

`GOOPSC_GITHUB_USER`

required: true

Github username

`GOOPSC_GITHUB_TOKEN`

required: true

Personal token for git repository with write access.

## Travis CI example
```yaml
after_success:
- goops nightly

deploy:
  provider: releases
  api_key: "$GOOPSC_GITHUB_TOKEN"
  file_glob: true
  file: bin/*
  prerelease: true
  overwrite: true
  skip_cleanup: true
  on:
    tags: true

```


