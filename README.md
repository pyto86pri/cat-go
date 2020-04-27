# cat go
cat command implemented with Golang.

## Build
```
$ go build
```

## Usage
Usage is same as cat command.
```
$ ./cat-go fileA fileB

aaaaa
aaaaa
bbbbb
bbbbb
```

### Options
#### -n
Display with line number.
```
$ ./cat-go -n fileA

1:aaaaa
2:aaaaa
```
#### -b
Display with line number skipping blank lines.
```
$ ./cat-go -b fileA

1:aaaaa

2:aaaaa
```
#### -s
Display lines squeezing continuous blank lines to one blank line.
```
$ ./cat-go -s fileA

aaaaa

aaaaa
```
#### -e
Display lines with "$" at the end of each lines.
```
$ ./cat-go -e fileA

aaaaa$
$
aaaaa$
```
#### -t
Display lines with tab replaced with "^I".
```
$ ./cat-go -t fileA

^Iaaaaa

aaaaa
```