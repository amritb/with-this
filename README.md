[![Build Status](https://travis-ci.org/amritb/with-this.svg?branch=master)](https://travis-ci.org/amritb/with-this)

# with-this

`with` let's you run any shell command with variables, in parallel.

Example:
```
$ with -v "$(cat myurls.txt)" "curl -L this"
```
