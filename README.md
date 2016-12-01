# ts

`ts` is the command line client for Timberslide. This client is under development
and may change at any time.

## Build Options

* [Download a binary](https://github.com/timberslide/ts/releases) 
* Use the Go package manager: `go get github.com/timberslide/ts`
* Build from source with `go install` (Requires Go to be installed). 

Tested under Go 1.7.

## Configuration

~/.timberslide/config

```
host = "gw.timberslide.com:443"
token = "<yourtokenhere>"
```

Contact info@timberslide.com for a token.

## Usage

Send stream

`echo "Hello Timberslide" | ts send hellots`

Receive a stream

`ts get hellots`
