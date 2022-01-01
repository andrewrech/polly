[![GoDoc](https://godoc.org/github.com/andrewrech/polly?status.svg)](https://godoc.org/github.com/andrewrech/polly) [![](https://goreportcard.com/badge/github.com/andrewrech/polly)](https://goreportcard.com/report/github.com/andrewrech/polly) ![](https://img.shields.io/badge/docker-andrewrech/polly:0.0.4-blue?style=plastic&logo=docker)

# polly

`polly` is a simple CLI tool to transform academic plain text files into audio using [AWS Polly](https://aws.amazon.com/polly/). Extraneous text such as references are removed.

## Installation

See [Releases](https://github.com/andrewrech/polly/releases).

## Usage

See `polly -h`:

```
Transform academic plain text files into audio using AWS Polly.

Usage of polly:

Defaults:
  -bucket string
        Output S3 bucket name (default "my-bucket")
  -dry-run
        Print TTS to file without processing. (default true)
  -engine string
        TTS engine (standard or neural) (default "neural")
  -format string
        Output format (mp3, ogg_vorbis, or pcm) (default "mp3")
  -input string
        Filename containing text to convert (default "input.txt")
  -prefix string
        Output S3 bucket prefix (default "<filename>")
  -voice string
        Voice to use for synthesis (Joanna, Salli, Kendra, Matthew, Amy [British], Brian [British], Olivia [Australian]) (default "Joanna")

Optional environmental variables:

    export AWS_SHARED_CREDENTIALS_PROFILE=default
    export AWS_SNS_TOPIC_ARN=my_topic_arn
```

## Authors

- [Andrew J. Rech](mailto:rech@rech.io)

## License

GNU Lesser General Public License v3.0
