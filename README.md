# Mainflux Auth Service

[![License](https://img.shields.io/badge/license-Apache%20v2.0-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.org/mainflux/mainflux-auth.svg?branch=master)](https://travis-ci.org/mainflux/mainflux-auth)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mainflux/mainflux-auth)](https://goreportcard.com/report/github.com/Mainflux/mainflux-auth)
[![Join the chat at https://gitter.im/Mainflux/mainflux](https://badges.gitter.im/Join%20Chat.svg)][gitter]

Mainflux IoT authentication and authorization management service.

### Install

Mainflux Auth Server uses [Redis](https://redis.io/), so insure that it is installed on your system.

Installing Mainflux Auth Server is trivial [`go get`](https://golang.org/cmd/go/):
```bash
go get github.com/mainflux/mainflux-auth
$GOBIN/mainflux-auth
```

If you are new to Go, more information about setting-up environment and fetching Mainflux code can be found [here](https://github.com/mainflux/mainflux-core-doc/blob/master/goenv.md).

### Community
#### Mailing lists
- [mainflux-dev][google-dev] - developers related. Discussions about development of Mainflux IoT itself.
- [mainflux-user][google-user] - general discussion and support. If you do not participate in development
    of Mainflux cloud infrastructure, this is probably what you're looking for.

#### IRC
[Mainflux Gitter][gitter]

#### Twitter
[@mainflux][twitter]

### License
[Apache License, version 2.0](LICENSE)

[wiki]: https://github.com/Mainflux/mainflux/wiki
[google-dev]: https://groups.google.com/forum/#!forum/mainflux-dev
[google-user]: https://groups.google.com/forum/#!forum/mainflux-user
[twitter]: https://twitter.com/mainflux
[gitter]: https://gitter.im/Mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge
