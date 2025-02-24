# ubuntu-autoinstall

[![Build, Release, and Attest](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/release.yaml/badge.svg)](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/release.yaml) [![CodeQL](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/github-code-scanning/codeql) [![Nightly Build and Publish](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/nightly.yaml/badge.svg)](https://github.com/jdfalk/ubuntu-autoinstall-webhook/actions/workflows/nightly.yaml)

A simple golang application to process ubuntu auto install reporting events

## Build

```shell
go build -o webhook
./webhook serve
```
