FROM golang:latest
COPY polly /
ENTRYPOINT ["/polly"]
