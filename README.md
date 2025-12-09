# gh-commit-analyzer

Interactive CLI to browse commits from your GitHub repositories with semantic search.

## Requirements

- [gh](https://cli.github.com/) CLI installed and authenticated
- [ck](https://github.com/BeaconBay/ck) (optional) for semantic search

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

### Optional: Semantic Search

Install `ck` for semantic commit message search:

```bash
cargo install ck-search
```

## Usage

```bash
gh-commit-analyzer
```

## Features

- Browse commits across multiple repositories
- Filter by date range, author, and branch
- Semantic search on commit messages (requires `ck`)
- Pagination with automatic and manual load more

## Semantic Search

Enter queries in the "Semantic" field to find commits by meaning:

- `bug fix` → finds commits about fixes, patches, corrections
- `refactoring` → finds cleanup, restructure, reorganize commits
- `performance` → finds optimization, speed improvements

## Controls

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate |
| `space` | Select |
| `enter` | Confirm/Expand |
| `tab` | Next field |
| `/` | Search |
| `n` | Load more commits |
| `r` | Restart |
| `d` | Default branch (all repos) |
| `q` | Quit |
