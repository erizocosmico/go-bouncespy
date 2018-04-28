# go-bouncespy [![Build Status](https://travis-ci.org/erizocosmico/go-bouncespy.svg?branch=master)](https://travis-ci.org/erizocosmico/go-bouncespy) [![GoDoc](https://godoc.org/gopkg.in/erizocosmico/go-bouncespy.v1?status.svg)](http://godoc.org/gopkg.in/erizocosmico/go-bouncespy.v1) [![codebeat badge](https://codebeat.co/badges/fe9d975f-de89-4c94-993e-8d9049833a0a)](https://codebeat.co/projects/github-com-erizocosmico-go-bouncespy)
Golang library to find the reason why your emails bounced according to the [RFC3463](https://tools.ietf.org/html/rfc3463#section-3) and [RFC821](https://tools.ietf.org/html/rfc821#section-4.2.2).


## Install

```
go get gopkg.in/erizocosmico/go-bouncespy.v1
```

## Usage

```go
import (
        "gopkg.in/erizocosmico/go-bouncespy.v1"
)

func main() {
        result := bouncespy.Analyze(emailHeaders, emailBody)
}
```
