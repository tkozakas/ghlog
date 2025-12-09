# gh-commit-analyzer

Interactive CLI to browse commits from your GitHub repositories.

## Requirements

- [gh](https://cli.github.com/) CLI installed and authenticated

## Installation

```bash
go install github.com/tkozakas/gh-commit-analyzer@latest
```

Or build from source:

```bash
git clone https://github.com/tkozakas/gh-commit-analyzer.git
cd gh-commit-analyzer
go build -o gh-commit-analyzer .
```

## Usage

```bash
gh-commit-analyzer
```

## Controls

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate |
| `space` | Select |
| `enter` | Confirm |
| `tab` | Next field |
| `/` | Search |
| `d` | Default branch (all repos) |
| `q` | Quit |
