# gitlab-cli
Gitlab REST API client.

## Install
### Linux
```bash
sudo curl -L "https://github.com/sotomskir/gitlab-cli/releases/download/0.1.0/gitlab-cli-$(uname -s)-$(uname -m)" -o /usr/local/bin/gitlab-cli
```
### Other platforms
You can download pre-build binary here: https://github.com/sotomskir/gitlab-cli/releases

## Usage
### Gitlab API auth
If You run gitlab-cli from Gitlab pipeline it will read API URL and auth token from env variables `CI_GITLAB_TOKEN` and `CI_API_V4_URL`.
You can also set these by `--token` and `--server` flags like this:
```bash
gitlab-cli --server https://gitlab.com/api/v4 --token XKSGhwe_dsbu27Gs
```
Another way is to use configuration file (`~/.gitlab-cli.yaml` by default)
```
# ~/.gitlab-cli.yaml
gitlab_token: XKSGhwe_dsbu27Gs
api_v4_url: https://gitlab.com/api/v4
```

### Commands
To list available commands type:
```
gitlab-cli --help
```

```
gitlab-cli --help
Gitlab REST API command line client.

Usage:
  gitlab-cli [command]

Available Commands:
  completion  Generates bash completion scripts
  help        Help about any command
  issues      List issue keys mentioned in merge request title and description
  version     Generate semantic version for current HEAD

Flags:
      --config string   config file (default is $HOME/.gitlab-cli.yaml)
  -h, --help            help for gitlab-cli
      --no-color        Disable ANSI color output
  -s, --server string   Gitlab API Url
  -t, --toggle          Help message for toggle
  -a, --token string    Gitlab API auth token

Use "gitlab-cli [command] --help" for more information about a command.
```
Each command has it's own help:
```
gitlab-cli version --help
```

```
gitlab-cli version --help
Generate semantic version for current HEAD.
Version generation is based on git tags.
If current HEAD is tagged then tag will be used as version.
Else command will lookup for previous tag bump it's minor version, reset patch version and append '-SNAPSHOT'
When there are no tags found version will be '0.1.0-SNAPSHOT'

Usage:
  gitlab-cli version [flags]

Aliases:
  version, v

Flags:
  -h, --help   help for version

Global Flags:
      --config string   config file (default is $HOME/.gitlab-cli.yaml)
      --no-color        Disable ANSI color output
```

## Bash completion
To load completion run
```bash
. <(gitlab-cli completion)
```

To configure your bash shell to load completions for each session add following line to your bashrc
 ~/.bashrc or ~/.profile
```bash
. <(gitlab-cli completion)
```

## Configuration
gitlab-cli read configuration from `~/.gitlab-cli.yaml` file by default. 
You can override this behaviour by `--config` flag. Most of the configuration options can be overridden by passing flag 
to gitlab-cli or by environment variable. Flags has highest precedence and configuration options has lowest precedence.
Table below shows all configuration options and corresponding flags and environment variables.

|Configuration option|Environment variable|Flag      |
|--------------------|--------------------|----------|
|api_v4_url          |CI_API_V4_URL       |--server  |
|gitlab_token        |CI_GITLAB_TOKEN     |--token   |
|merge_request_iid   |CI_MERGE_REQUEST_IID|--mr      |
|project_id          |CI_PROJECT_ID       |--project |
