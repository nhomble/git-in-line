# git-in-line

sort list of commit hashes by their revision order in a git repository.

## usage

```bash
gil <repo-path> <commit-file> [flags]
```

### examples

```bash
gil /path/to/repo commits.txt
```

## build

```bash
go build -o gil ./src
```

## test

```bash
./integration-tests/run-tests.sh
```