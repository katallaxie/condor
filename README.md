# :partying_face: Template Go

> This is the standard template for Go projects of @katallaxie.

## Get Started

```bash
make setup MODULE_NAME=github.com/username/repo
```

Features

* [Developer Containers](https://code.visualstudio.com/docs/remote/containers)
* [Editorconfig](https://editorconfig.org)
* [GoReleaser](https://goreleaser.com)

> You can `sh scripts/postCreateCommand.sh` if you are not running in a remote container or on [Codespaces](https://github.com/features/codespaces).

## Development

This template supports `Makefile` and `run` build tooling.

```bash
# build `main.go`
make 
```

Other available targets are

* `build`
* `fmt`
* `lint`
* `vet`
* `generate`
* `clean`

The convention is to use `make` to run the build.

Happy coding!
