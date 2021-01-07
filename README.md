[![GoDoc](https://godoc.org/github.com/andrewrech/polly?status.svg)](https://godoc.org/github.com/andrewrech/polly) ![](https://img.shields.io/badge/version-0.0.2-blue.svg) [![](https://goreportcard.com/badge/github.com/andrewrech/polly)](https://goreportcard.com/report/github.com/andrewrech/polly)

# polly

`polly` is a simple CLI tool to transform academic plain text files into audio using [AWS Polly](https://aws.amazon.com/polly/). Extraneous text such as references are removed.

## Installation

```zsh
go get -u -v github.com/andrewrech/polly
```

## Usage

See `polly -h`.

## Documentation

See [docs](https://github.com/andrewrech/polly/blob/main/docs.md).

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
