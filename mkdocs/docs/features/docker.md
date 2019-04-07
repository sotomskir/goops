## Description
Automatic docker stable & latest tags based on selected semver strategy.

## Output variables
```console
none
```

## Configuration defaults
```console
GOOPSC_DOCKER=false
GOOPSC_SEMVER_STRATEGY=gitlab-flow
```

## gitlab-flow strategy
All tagged builds are stable.
All builds from master are latest.

## Usage
```console
$ goops docker build -t $DOCKER_IMAGE:$GOOPS_SEMVER .
$ goops docker push $DOCKER_IMAGE:$GOOPS_SEMVER
```
