# Kandinsky

[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![Build Status](https://travis-ci.org/andersnormal/kandinsky.svg?branch=master)](https://travis-ci.org/andersnormal/kandinsky)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


A relay between WebSocket protocol and TCP services. Thus, it allows to use the [WebSocket Protocol](https://tools.ietf.org/html/rfc6455) and make a connection to the relay, which copies all TCP traffic from and to the relayed service.

## Usage

> use `--help` to see all available options

```bash
Usage:
  kandinsky [flags]

Flags:
      --addr string    address to listen on (default ":8080")
  -h, --help           help for kandinsky
      --relay string   address to relay (e.g. nats:4222
      --verbose        enable verbose output
```

## License
[Apache 2.0](/LICENSE)
