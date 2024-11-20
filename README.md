# Git File Downloader (gdown)

![Go](https://github.com/haevg-rz/git-file-downloader/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/haevg-rz/git-file-downloader/branch/master/graph/badge.svg)](https://codecov.io/gh/haevg-rz/git-file-downloader)
[![Go Report Card](https://goreportcard.com/badge/github.com/haevg-rz/git-file-downloader)](https://goreportcard.com/report/github.com/haevg-rz/git-file-downloader)

Download a file or folder from a git hosting service and save it to disk if file is different, to ensure that the configuration files (or other files) on your servers are always up to date.

## Latest

See Releases

## Using

```plain
PS C:\DevGit\git-file-downloader> .\gdown.exe --help
git-file-downloader

Usage:
  gdown <github|gitlab|azure|help|version> [flags]
  gdown [command]

Available Commands:
  azure       retrieves data from azure dev ops
  completion  Generate the autocompletion script for the specified shell
  github      retrieves data from github
  gitlab      retrieves data from gitlab
  help        Help about any command

Flags:
      --branch string        Branch name (default "main")
      --exclude string       Exclude this regex pattern
  -h, --help                 help for gdown
      --include string       Include this regex pattern
      --logfile              Write to file instead of stdout
      --out string           Path to write file to disk
      --remote-path string   Path to file/folder from remote source
      --token string         Private-Token with access right for "api" and "read_repository", role must be minimum "Reporter"
      --url string           url to Api v4, like https://my-git-lab-server.local/api/v4/
      --user-agent string    User agent (default "Go-http-client/1.1")
  -v, --verbosity int        Set verbosity level (0-3) (default 3)

Use "gdown [command] --help" for more information about a command.
```

## Use Case

### Download files from your gitlab repository

You want to have the benefits from git to manage your config files, but don't want git installed on your system?
With this (windows and linux) tool you can now download these config files from an on-promise instance of your choice (GitHub, AzureGit, GitLab) and save them to your disk.

You can download folders/files from different repositories and providers, and save them in the output path of your choice.

Files will **only** be replaced if the hashes are different (from disk to git).

## Using
The tool is divided into the following subcommands:
- `github`
- `gitlab`
- `azure`
- `public` (soon)

Depending on which git-provider you choose, you will need to provide different mandatory flags.
For more info you can use the `help` command on any of the subcommands. This will list all the global and local flags for the command with a short description.
For example:
```plain
PS C:\DevGit\git-file-downloader> go run .\cmd\gdown\main.go help github
retrieves data from github

Usage:
gdown github <file|folder> <flags> [flags]

Flags:
-h, --help           help for github
--owner string   repo owner
--repo string    repo name

Global Flags:
--branch string        Branch name (default "main")
--exclude string       Exclude this regex pattern
--include string       Include this regex pattern
--logfile              Write to file instead of stdout
--out string           Path to write file to disk
--remote-path string   Path to file/folder from remote source
--token string         Private-Token with access right for "api" and "read_repository", role must be minimum "Reporter"
--url string           url to Api v4, like https://my-git-lab-server.local/api/v4/
--user-agent string    User agent (default "Go-http-client/1.1")
-v, --verbosity int        Set verbosity level (0-3) (default 3)
2024/11/20 16:36:24 exit code: 0
```

After providing the subcommand, you will need to provide an argument specifying whether you want to download a single file or a folder.

**Important**

If you choose to download a folder, be aware that gdown will recursively download any possible nested folders aswell. (will soon be an optional flag)

## Example
### Download folder from your gitlab repository
**Working example!**

```bat
.\gdown.exe gitlab folder --token AZURE_PAT --out ./clone-gitlab --remote-path / --project 64241402
```

```log
2024/11/20 16:14:14 Sync 2 files, from remote folder /
2024/11/20 16:14:15 Sync 3 files, from remote folder config
2024/11/20 16:14:15 Sync 3 files, from remote folder config/foo
2024/11/20 16:14:15 Sync 1 files, from remote folder config/foo/nested
2024/11/20 16:14:16 Created File: 'clone-gitlab/config/foo/nested/.gitkeep' because it didn't exist
2024/11/20 16:14:16 Wrote file: config/foo/nested/.gitkeep because is new or updated
2024/11/20 16:14:16 Created File: 'clone-gitlab/config/foo/.gitkeep' because it didn't exist
2024/11/20 16:14:16 Wrote file: config/foo/.gitkeep because is new or updated
2024/11/20 16:14:16 Created File: 'clone-gitlab/config/foo/bar.txt' because it didn't exist
2024/11/20 16:14:16 Wrote file: config/foo/bar.txt because is new or updated
2024/11/20 16:14:16 Created File: 'clone-gitlab/config/.gitkeep' because it didn't exist
2024/11/20 16:14:16 Wrote file: config/.gitkeep because is new or updated
2024/11/20 16:14:17 Created File: 'clone-gitlab/config/greetings.txt' because it didn't exist
2024/11/20 16:14:17 Wrote file: config/greetings.txt because is new or updated
2024/11/20 16:14:17 Created File: 'clone-gitlab/README.md' because it didn't exist
2024/11/20 16:14:17 Wrote file: README.md because is new or updated
2024/11/20 16:14:17 synced file(s) successfully
2024/11/20 16:14:17 exit code: 0
```

### Download file from your gitlab repository
**Working example!**

See https://gitlab.com/mdriessen/downloader-test for local testing.

```bat
.\gdown.exe gitlab file --token AZURE_PAT --out ./greetings.txt --remote-path config/greetings.txt --project 64241402
```

```log
2024/11/20 16:12:11 Created File: './greetings.txt' because it didn't exist
2024/11/20 16:12:11 synced file(s) successfully
2024/11/20 16:12:11 exit code: 0
```

## Contributing

- Github Copilot [.github/copilot-instructions.md](.github/copilot-instructions.md)

## Technical

### Scopes for personal access tokens

#### Gitlab

- `read_repository`: Allows read-access to the repository files.
- `api`: Allows read-write access to the repository files.

#### Github

- `Contents`: minimum `Read-Only`

#### Azure

- `Code`: minimum `Read`

### TLS Security

**Will be a config switch, soon.**

```go
tr := &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
```
