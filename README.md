[![GoDoc](https://godoc.org/github.com/andrewrech/polly?status.svg)](https://godoc.org/github.com/andrewrech/polly) [![](https://goreportcard.com/badge/github.com/andrewrech/polly)](https://goreportcard.com/report/github.com/andrewrech/polly) ![](https://img.shields.io/badge/docker-andrewrech/polly:0.0.4-blue?style=plastic&logo=docker)

# polly

`polly` is a simple CLI tool to transform academic plain text files into audio using [AWS Polly](https://aws.amazon.com/polly/). Extraneous text such as references are removed.

## Installation

See [Releases](https://github.com/andrewrech/polly/releases).

```zsh
go get -u -v github.com/andrewrech/polly
```

## Usage

See `polly -h` or [documentation](https://github.com/andrewrech/polly/blob/main/docs.md)).

## Testing

```zsh
git clone https://github.com/andrewrech/polly &&
cd polly

go test
```

## Authors

- [Andrew J. Rech](mailto:rech@rech.io)

## License

GNU Lesser General Public License v3.0
