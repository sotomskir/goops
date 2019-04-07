## Gitlab CI
```yaml
image: docker:stable
variables:
  # When using dind service we need to instruct docker, to talk with the
  # daemon started inside of the service. The daemon is available with
  # a network connection instead of the default /var/run/docker.sock socket.
  #
  # The 'docker' hostname is the alias of the service container as described at
  # https://docs.gitlab.com/ee/ci/docker/using_docker_images.html#accessing-the-services
  #
  # Note that if you're using Kubernetes executor, the variable should be set to
  # tcp://localhost:2375 because of how Kubernetes executor connects services
  # to the job container
  DOCKER_HOST: tcp://docker:2375/
  # When using dind, it's wise to use the overlayfs driver for
  # improved performance.
  DOCKER_DRIVER: overlay2

services:
- docker:dind

stages:
- common
- build
- deploy
- jira

cache:
  key: "$CI_PIPELINE_ID"
  paths:
  # keep env file between stages
  - ./gitlab.env

# Generate semantic version for current build,
common:
  stage: common
  image: sotomski/gitlab-tools
  script:
  - goops pipeline common
  only:
  - merge_requests
  - master
  - tags
  - /^.*-stable$/

lint:
  stage: build
  image: node:10
  script:
  - yarn install
  - node_modules/.bin/ng lint
  only:
  - merge_requests
  - master
  - tags
  - /^.*-stable$/

test:
  stage: build
  image: sotomski/node:10-chrome
  script:
  - yarn install
  - node_modules/.bin/ng test --browsers ChromeHeadlessNoSandbox --source-map=false --watch=false
  only:
  - merge_requests
  - master
  - tags
  - /^.*-stable$/

# build
build-push:
  stage: build
  image: sotomski/gitlab-tools:dind
  script:
  - source gitlab.env
  - docker login -u gitlab-ci-token -p $CI_JOB_TOKEN registry.gitlab.com
  - goops pipeline docker build -t $CI_REGISTRY_IMAGE:$CI_SEMVER .
  - goops pipeline docker push $CI_REGISTRY_IMAGE:$CI_SEMVER
  only:
  - merge_requests
  - master
  - tags
  - /^.*-stable$/

jira-cr:
  stage: jira
  image: sotomski/gitlab-tools
  script:
  - source gitlab.env
  - goops pipeline jira transition "code review"
  only:
  - merge_requests

jira-in-test:
  stage: jira
  image: sotomski/gitlab-tools
  script:
  - source gitlab.env
  - goops pipeline jira transition "in test"
  only:
  - master

# release Jira version
release:
  stage: deploy
  image: sotomski/gitlab-tools
  script:
  - source gitlab.env
  - jira-cli version release $JIRA_PROJECT_KEY $CI_SEMVER_RELEASE
  only:
  - tags

# deploy staging environment
deploy_staging:
  stage: deploy
  image: alpine
  script:
  - source gitlab.env
  - echo deploy version $CI_SEMVER to staging environment
  environment:
    name: staging
    url: $STAGING_URL
  only:
  - master

# deploy pre-prod environment
deploy_pre-prod:
  stage: deploy
  image: alpine
  script:
  - source gitlab.env
  - echo deploy version $CI_SEMVER to pre-prod environment
  environment:
    name: pre-prod
    url: $PRE-PROD_URL
  when: manual
  only:
  - master
  - tags
  - /^.*-stable$/

# deploy prod environment
deploy_prod:
  stage: deploy
  image: alpine
  script:
  - source gitlab.env
  - echo deploy version $CI_SEMVER to prod
  environment:
    name: pre-prod
    url: $PROD_URL
  when: manual
  only:
  - tags
```

## Jenkins
## Travis