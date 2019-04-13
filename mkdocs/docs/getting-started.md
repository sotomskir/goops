## Installation

**Linux and Mac**

To install latest nightly build of goops run below command. 
If you wish to install stable version please replace `nightly` with specific version from 
[releases page](https://github.com/sotomskir/goops/releases).

```console
$ sudo curl -L "https://github.com/sotomskir/goops/releases/download/nightly/goops-$(uname -s)-$(uname -m)" -o /usr/local/bin/goops && chmod +x /usr/local/bin/goops 
```

**Windows**

You can download pre build exe file from 
[releases page](https://github.com/sotomskir/goops/releases).

## CI pipeline configuration
goops is working on environment variables. 
To initialize goops variables setenv command should be run on the beginning of pipeline. 
```console
$ . <(goops setenv)
```
This command will set goops environment variables.

!!! note
    goops variables should be persisted between pipeline stages.
    Many CI/CD tools will reset variables on each stage. It's up to you
    to make sure that variables are persisted between stages.

## Persisting variables between pipeline stages
`goops setenv` command will also save variables to `.goops.env` file. 
To restore variables from file run:
```console
$ source .goops.env
```
You must ensure that `.goops.env` file is persisted between stages. Below you will find couple
of examples how to achieve this in various CI tools.

**Gitlab CI**
`.gitlab-ci.yml`
```yaml
cache:
  key: "$CI_PIPELINE_ID"
  paths:
  # keep env file between stages
  - ./.goops.env
```
TODO: add travis example

## Configuration methods
goops has two methods of configuration

* environment variables
* configuration file stored in repo `.goops.yml`

!!! note
    Environment variables should be uppercase e.g. GOOPS_CI_TYPE=jenkins while 
    configuration file variables should be lowercase e.g. goops_ci_type: jenkins

## Gitlab environment variables
Environment variables in Gitlab can be configured on group or project level settings > CI/CD
![Gitlab variables](./img/gitlab_env_variables.png?raw=true "Gitlab variables")

## Travis CI environment variables
Go to project settings > environment variables section
![Travis variables](./img/travis_env.png?raw=true "Travis variables")

## Jenkins environment variables
Manage Jenkins > Configure System > Global variables
![Jenkins variables](./img/jenkins_env.png?raw=true "Jenkins variables")
