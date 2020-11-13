package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/google/go-cmp/cmp"
)

func TestTTS(t *testing.T) {
	vars, err := loadVars()
	if err != nil {
		log.Fatalln(err)
	}

	fileName := "./testdata/test.txt"
	engine := "standard"
	format := "mp3"
	outputS3BucketName := "rech-public"
	outputS3BucketPrefix := getFnPrefix(&fileName)
	voiceID := "Joanna"

	// Open text file and get it's contents as a string
	text, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Got error opening file:", err.Error())
	}

	// Convert bytes to string
	s := string(text)
	s = TTSformat(s)

	input := polly.StartSpeechSynthesisTaskInput{
		Engine:             aws.String(engine),
		OutputFormat:       aws.String(format),
		OutputS3BucketName: aws.String(outputS3BucketName),
		OutputS3KeyPrefix:  outputS3BucketPrefix,
		SnsTopicArn:        aws.String(vars.snsTopic),
		Text:               aws.String(s),
		VoiceId:            aws.String(voiceID),
	}

	c, output := getInput(input, vars)

	t.Run("Task output generation", func(t *testing.T) {
		diff := cmp.Diff(getTaskStatus(c, output), "scheduled")
		if diff != "" {
			t.Fatalf(diff)
		}
	})
}

func TestTTSformat(t *testing.T) {
	s := `Title

Summary
Some text (NCT2345432) more (123) text (Author et al, 2020), figure (Fig. 3).22,23`

	sOut := TTSformat(s)

	t.Run("Task output generation", func(t *testing.T) {
		diff := cmp.Diff(sOut, "Title... Next. Summary\nSome text more text, figure.")
		if diff != "" {
			t.Fatalf(diff)
		}
	})
	log.Println()
}
