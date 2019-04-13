## Description

When enabled setenv command will fetch jira issue keys from merge request title, description and commit messages linked to merge request.
Each issue found will be linked to jira version. If Jira version does not exists it will be created, along with deployment issue.
`GOOPS_SEMVER_RELEASE` variable must be set, by semver feature or manually. 

## Output variables

```console
GOOPS_JIRA_ISSUES=TEST-1 TEST-2 TEST-3
```

## Configuration defaults

```console
GOOPSC_JIRA=false
GOOPSC_JIRA_PROJECT_KEY=
GOOPSC_JIRA_SERVER_URL=
GOOPSC_JIRA_USER=
GOOPSC_JIRA_PASSWORD=
GOOPSC_JIRA_VERSION_ASSIGN=true
GOOPSC_JIRA_VERSION_CREATE=true
GOOPSC_JIRA_CREATE_DEPLOYMENT_ISSUE=true
GOOPSC_JIRA_ISSUE_TRANSITION=true
GOOPSC_JIRA_WORKFLOW=workflow.yaml
GOOPSC_JIRA_WORKFLOW_CONTENT=
GOOPSC_JIRA_STRATEGY=gerrit
```
`GOOPSC_JIRA`

Enable Jira integration

`GOOPSC_JIRA_PROJECT_KEY`

required: true

Jira project key

`GOOPSC_JIRA_SERVER_URL`

required: true

Jira server url e.g. https://jira.example.com

`GOOPSC_JIRA_USER`

required: true

Jira username

`GOOPSC_JIRA_PASSWORD`

required: true

Jira password

`GOOPSC_JIRA_VERSION_ASSIGN` 

Assign issues to Jira version. Semver feature must be enabled or GOOPS_SEMVER_RELEASE variable set.
Issue list is taken from GOOPS_JIRA_ISSUES variable.

`GOOPSC_JIRA_VERSION_CREATE`

Create Jira version. When assigning issue to version and version doesn't exists it will be created.

`GOOPSC_JIRA_CREATE_DEPLOYMENT_ISSUE`

When creating Jira version create deployment issue assigned to version.

`GOOPSC_JIRA_ISSUE_TRANSITION`

Transition Jira issues.

`GOOPSC_JIRA_WORKFLOW`

Path to workflow definition. Can be local file or remote http path.

## Transitioning issues

```console
$ goops jira transition 'target state'
```

jira transition command require workflow definition in yaml file. 
Default filename is `workflow.yaml` and can be overridden by --workflow flag or `GOOPSC_JIRA_WORKFLOW` variable.
Remote http url is also accepted.

**workflow structure**

```yaml
workflow:
  source status:
    target status: transition name
    default: default transition name
```

**example workflow definition**

```yaml
workflow:
  code review:
    default: ready to test
  in test:
    done: done
    default: bug found
  to do:
    rejected: reject
    default: start progress
  in progress:
    default: code review
  done:
    default: reopen
  rejected:
    default: reopen
```

**corresponding Jira workflow**

![Alt text](../img/workflow.png?raw=true "Example Jira workflow")

**Workflow from env variable**

Alternatively workflow file content can be passed by `GOOPSC_JIRA_WORKFLOW_CONTENT` environment variable.
```yaml
export JIRA_WORKFLOW_CONTENT=$(cat  <<- EOM
workflow:
  code review:
    default: ready to test
  in test:
    done: done
    default: bug found
  to do:
    rejected: reject
    default: start progress
  in progress:
    default: code review
  done:
    default: reopen
  rejected:
    default: reopen
EOM
)
```
