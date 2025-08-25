# git-in-line

sort list of commit hashes by their revision order in a git repository.

## Usage

```bash
gil <repo-path> <commit-file> [flags]
```

### examples

```bash
gil /path/to/repo commits.txt
```

## Build

```bash
go build -o gil ./src
```

## Test

```bash
./integration-tests/run-tests.sh
```