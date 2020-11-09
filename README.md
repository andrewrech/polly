[![GoDoc](https://godoc.org/github.com/andrewrech/polly?status.svg)](https://godoc.org/github.com/andrewrech/polly) ![](https://img.shields.io/badge/version-0.0.2-blue.svg) [![](https://goreportcard.com/badge/github.com/andrewrech/polly)](https://goreportcard.com/report/github.com/andrewrech/polly)

# polly

`polly` is a simple CLI tool to transform academic plain text files into audio using [AWS Polly](https://aws.amazon.com/polly/). Extraneous text such as references are removed.

## Installation

```zsh
go get -u -v github.com/andrewrech/polly
```

## Usage

```zsh
polly -h

```

```
Defaults:
  -bucket string
        Output S3 bucket name (default "my-bucket")
  -dry-run
        Print TTS text without uploading?
  -engine string
        TTS engine (standard or neural) (default "neural")
  -format string
        Output format (mp3, ogg_vorbis, or pcm) (default "mp3")
  -input string
        Filename containing text to convert (default "input.txt")
  -prefix string
        Output S3 bucket prefix (default "<filename>")
  -voice string
        Voice to use for synthesis (Joanna, Salli, Kendra, Matthew (default "Joanna")

Environmental variables:

    export AWS_ACCESS_KEY_ID=my_iam_access_key
    export AWS_SECRET_ACCESS_KEY=my_iam_secret
    export AWS_SNS_TOPIC_ARN=my_topic_arn
    export AWS_DEFAULT_REGION=my_region_name
    export AWS_SESSION_TOKEN=my_iam_session_token [optional]
```

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
