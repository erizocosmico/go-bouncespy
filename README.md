# go-bouncespy [![Build Status](https://travis-ci.org/mvader/go-bouncespy.svg?branch=master)](https://travis-ci.org/mvader/go-imapreader) [![GoDoc](https://godoc.org/gopkg.in/mvader/go-bouncespy.v1?status.svg)](http://godoc.org/gopkg.in/mvader/go-bouncespy.v1)
Golang library to find the reason why your emails bounced

## Install

```
go get gopkg.in/mvader/go-bouncespy.v1
```

## Usage

```go
import (
        "gopkg.in/mvader/go-bouncespy.v1"
)

func main() {
        result := bouncespy.Analyze(emailHeaders, emailBody)
}
```
