package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {

	exitVal := m.Run()

	// delete test output if it exists
	testOutput := []string{
		"www.google.com",
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	// delete downloaded test mp3
	for _, f := range files {
		fn := f.Name()

		if strings.HasPrefix(fn, "-testdata-test") && strings.HasSuffix(fn, ".mp3") {
			testOutput = append(testOutput, fn)
		}
	}

	for _, f := range testOutput {
		_, err := os.Stat(f)

		if err == nil {
			err := os.Remove(f)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	os.Exit(exitVal)
}

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

	t.Run("Task output handler", func(t *testing.T) {
		outputHandler(c, output)
	})

}

func TestDownload(t *testing.T) {

	download("https://www.google.com")

	_, err := os.Stat("www.google.com")

	fi, err := os.Stat("www.google.com")

	t.Run("download", func(t *testing.T) {
		if err != nil {
			t.Fatalf("Download file does not exist")
		}

		if fi.Size() == 0 {
			t.Fatalf("Download file empty")
		}
	})

}

func TestTTSformat(t *testing.T) {
	s := `Title

Summary
Some text (NCT2345432) moré (123) text (Author et al, 2020), figure (Fig. 3).22,23 Tumor cells2.`

	sOut := TTSformat(s)

	t.Run("Task output generation", func(t *testing.T) {
		diff := cmp.Diff(sOut, "Title\n\nSummary\nSome text more text, figure.22,23 Tumor cells2.\n")
		if diff != "" {
			t.Fatalf(diff)
		}
	})
	log.Println()
}
