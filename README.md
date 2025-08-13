# gitlabctl

A command-line tool for managing GitLab projects and branches.

## Features

- Verify branch existence across multiple GitLab projects
- Create branches across multiple GitLab projects
- Bulk operations support

## Installation

```bash
go install github.com/gregaf/gitlabctl/cmd/cli@latest
```

## Configuration

The tool requires a configuration file with your GitLab access token. Create a JSON configuration file:

```json
{
  "gitlabUrl": "https://gitlab.example.com",
  "accessToken": "your-gitlab-token"
}
```

## Usage

### Verify Branch Existence

```bash
gitlabctl verify --branch branch-name project1 project2
```

### Create Branch

```bash
gitlabctl create --branch new-branch --base main project1 project2
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.