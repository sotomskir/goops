## Description

Auto semantic version based on selected strategy.
Available strategies:

* github-flow
* gitlab-flow
* git-flow-branch

## Output variables
```console
GOOPS_SEMVER=1.2.3-SNAPSHOT
GOOPS_SEMVER_RELEASE=1.2.3
GOOPS_SEMVER_MAJOR=1
GOOPS_SEMVER_MINOR=2
GOOPS_SEMVER_PATCH=3
```

## Configuration defaults

```console
GOOPSC_SEMVER=false
GOOPSC_SEMVER_STRATEGY=github-flow
```

## gitlab-flow strategy

This strategy is designed for Gitlab flow with release branches. 
Which is described in [Gitlab documentation](https://docs.gitlab.com/ee/workflow/gitlab_flow.html#release-branches-with-gitlab-flow)

rules for master branch

1. Find previous tag. If there are no tags previous tag will be assumed as 0.0.0
2. Bump previous tag minor version and set patch version to 0.
3. If *-stable branch matching version exists bump minor version once more.
4. Append "-SNAPSHOT" to version.

rules for *-stable branches

1. If HEAD is tagged use tag as version.
2. Else find previous tag and bump patch version.
3. If tag not exists take version from branch name and set patch to 0.

| current branch | tag     | previousTag | stableBranch | version         | release version |
| -------------- |---------|-------------|--------------|-----------------|-----------------|
| master         |         |             |              | 0.1.0-SNAPSHOT  | 0.1.0           |
| master         |         |             | 0.1-stable   | 0.2.0-SNAPSHOT  | 0.2.0           |
| master         |         | 0.1.0       | 0.1-stable   | 0.2.0-SNAPSHOT  | 0.2.0           |
| master         |         | 0.1.1       | 0.1-stable   | 0.2.0-SNAPSHOT  | 0.2.0           |
| master         | 0.1.1   | 0.1.0       | 0.1-stable   | 0.2.0-SNAPSHOT  | 0.2.0           |
| 0.1-stable     |         |             | 0.1-stable   | 0.1.0-SNAPSHOT  | 0.1.0           |
| 0.1-stable     | 0.1.0   |             | 0.1-stable   | 0.1.0           | 0.1.0           |
| 0.1-stable     |         | 0.1.0       | 0.1-stable   | 0.1.1-SNAPSHOT  | 0.1.1           |
| 0.1-stable     | 0.1.1   | 0.1.0       | 0.1-stable   | 0.1.1           | 0.1.1           |
| 0.2-stable     |         |             | 0.2-stable   | 0.2.0-SNAPSHOT  | 0.2.0           |
