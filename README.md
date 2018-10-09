# with-this [![Build Status](https://travis-ci.org/amritb/with-this.svg?branch=master)](https://travis-ci.org/amritb/with-this)

`With` is a CLI to execute any shell command with multiple variables concurrently. xargs is a similar utility. Place the **this** keyword anywhere in your command and `with` will iterate through all the input values, replace **this** with a value and execute the resulting command.

If you have multiple **this** in your command, each one will be replaced (*todo: this should be controllable in a future release using flags*).

In case of errors, `with` will only report and continue with the remaining iterations/values.

## Installation
### macOS
```
$ brew install amritb/tap/with-this
```
### Linux / Windows
You can build it from source using `go build` and put it in your `$PATH`.
```
$ go build main.go -o with
```

## Usage
```
with --values "$(INPUT)" "COMMAND"

Flags:
  -h, --help            help for with
  -v, --values string   Iterate with these values
      --version         version for with

```

## Examples

You have a list of URLs in a text file and want to `curl` all of them in parallel, with one command:
```
$ with -v "$(cat myurls.txt)" "curl -L this"
```

You want to quickly check AWS instance status for all the *regions*:
```
$ with -v "$(cat myregions.txt)" "aws --region=this ec2 describe-instance-status"
```

You have a directory with lots of *kubeconfig* files and want to get pods from all the different clusters using `kubectl` command:
```
$ with -v "$(ls | grep config)" "kubectl --kubeconfig=this get pods"
```