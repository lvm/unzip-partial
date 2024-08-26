# unzip-partial

These tiny apps solve an insane issue i've found on macOS: you can't extract partially zip files. 

## Usage

### `unzip-partial`

```
$ unzip-partial
Usage: unzip-partial -zip <zip-file> -pattern <pattern> -output <output-dir>
```

### `unzip-partial-ls`

```
$ unzip-partial-ls
Usage: unzip-partial-ls -zip <zip-file>
```


## Build

```
go build -ldflags "-s -w" -o bin/unzip-partial cmd/unzip-partial/main.go
go build -ldflags "-s -w" -o bin/unzip-partial-ls cmd/unzip-partial-ls/main.go
```

## LICENSE 
See [LICENSE](LICENSE)
