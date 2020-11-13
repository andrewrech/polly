package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// usage prints package usage.
func usage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nTransform academic plain text files into audio using AWS Polly.\n")
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDefaults:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Environmental variables:

    export AWS_ACCESS_KEY_ID=my_iam_access_key
    export AWS_SECRET_ACCESS_KEY=my_iam_secret
    export AWS_SNS_TOPIC_ARN=my_topic_arn
    export AWS_DEFAULT_REGION=my_region_name
    export AWS_SESSION_TOKEN=my_iam_session_token [optional]

`)
	}
}

func main() {
	usage()

	// parse command line options
	fileName := flag.String("input", "input.txt", "Filename containing text to convert")
	engine := flag.String("engine", "neural", "TTS engine (standard or neural)")
	format := flag.String("format", "mp3", "Output format (mp3, ogg_vorbis, or pcm)")
	outputS3BucketName := flag.String("bucket", "my-bucket", "Output S3 bucket name")
	outputS3BucketPrefix := flag.String("prefix", "<filename>", "Output S3 bucket prefix")
	voiceID := flag.String("voice", "Joanna", "Voice to use for synthesis (Joanna, Salli, Kendra, Matthew, Amy [British], Brian [British], Olivia [Australian])")
	dryrun := flag.Bool("dry-run", false, "Print TTS text without uploading?")

	flag.Parse()

	vars, err := loadVars()
	if err != nil {
		log.Fatalln(err)
	}

	if *outputS3BucketPrefix == "<filename>" {
		outputS3BucketPrefix = getFnPrefix(fileName)
	}

	// Open text file and get it's contents as a string
	text, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Fatalln("Got error opening file:", err.Error())
	}

	// Convert bytes to string
	s := string(text)
	sOut := TTSformat(s)

	input := polly.StartSpeechSynthesisTaskInput{
		Engine:             engine,
		OutputFormat:       format,
		OutputS3BucketName: outputS3BucketName,
		OutputS3KeyPrefix:  outputS3BucketPrefix,
		SnsTopicArn:        aws.String(vars.snsTopic),
		Text:               aws.String(sOut),
		VoiceId:            voiceID,
	}

	// print text transformation without uploading
	if *dryrun {
		getDiff(s, sOut)
	}

	if !*dryrun {
		c, output := getInput(input, vars)

		log.Println(*output)

		outputHandler(c, output)
	}
}

// getDiff generates a text diff after text substitutions for better listening comprehension.
func getDiff(s, sOut string) {
	log.Println(sOut)

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(s, sOut, false)

	var buff bytes.Buffer

	// print insertions and deletions
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("\x1b[32m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("\x1b[31m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		}
	}

	log.Print(buff.String())
}

// getInput generates an AWS Polly task input.
func getInput(i polly.StartSpeechSynthesisTaskInput, vars envVars) (c *polly.Polly, output *polly.StartSpeechSynthesisTaskOutput) {
	// Initialize a session that the SDK uses to load
	config := aws.Config{
		Region:      aws.String(vars.region),
		Credentials: credentials.NewStaticCredentials(vars.id, vars.secret, vars.token),
	}

	// create S3 upload manager
	sess := session.Must(session.NewSession(&config))

	// Create Polly client
	c = polly.New(sess)

	output, err := c.StartSpeechSynthesisTask(&i)
	if err != nil {
		log.Fatalln("Got error calling SynthesizeSpeech:", err.Error())
	}

	return c, output
}

// getFnPrefix returns an AWS Polly audio file prefix based on the input text filename.
func getFnPrefix(fileName *string) (ret *string) {
	fNPrefix := strings.TrimSuffix(*fileName, path.Ext(*fileName))

	rmSpec := regexp.MustCompile("[^A-Za-z0-9]+")

	fNPrefix = rmSpec.ReplaceAllString(fNPrefix, "-")

	return (&fNPrefix)
}

// outputHandler waits for AWS Polly task completion and handles output, generating a download link.
func outputHandler(c *polly.Polly, output *polly.StartSpeechSynthesisTaskOutput) {
	re := regexp.MustCompile("scheduled|inProgress")

	// behavior is based on task status
	for re.MatchString(getTaskStatus(c, output)) {
		log.Println("working")
		time.Sleep(1 * time.Second)
	}

	if getTaskStatus(c, output) == "failed" {
		log.Println(*output.SynthesisTask.TaskStatus)
		log.Fatalln("Got error during processing SynthesizeSpeech:", output)
	}

	// notify and download the audio file when task is complete
	if getTaskStatus(c, output) == "completed" {
		log.Println(*output.SynthesisTask.TaskStatus)
		download(*output.SynthesisTask.OutputUri)
	}
}

// getTaskStatus retrieves the status of the AWS Polly synthesis task.
func getTaskStatus(c *polly.Polly, output *polly.StartSpeechSynthesisTaskOutput) (status string) {
	input := polly.GetSpeechSynthesisTaskInput{
		TaskId: output.SynthesisTask.TaskId,
	}

	ret, err := c.GetSpeechSynthesisTask(&input)
	if err != nil {
		log.Fatalln(err)
	}

	return *ret.SynthesisTask.TaskStatus
}

// download download a url to a local file.
func download(url string) {
	r, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	f := path.Base(url)

	// Create the file
	out, err := os.Create(f)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fi, err := out.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Downloaded", fi.Size(), "bytes")
}

type envVars struct {
	id       string
	secret   string
	snsTopic string
	region   string
	token    string
}

// loadVars loads required environmental variables.
func loadVars() (vars envVars, err error) {
	id, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		log.Fatalln()
		return vars, errors.New("AWS_ACCESS_KEY_ID is unset")
	}

	secret, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return vars, errors.New("AWS_SECRET_ACCESS_KEY is unset")
	}

	snsTopic, ok := os.LookupEnv("AWS_SNS_TOPIC_ARN")
	if !ok {
		return vars, errors.New("AWS_SNS_TOPIC_ARN is unset")
	}

	region, ok := os.LookupEnv("AWS_DEFAULT_REGION")
	if !ok {
		return vars, errors.New("AWS_DEFAULT_REGION is unset")
	}

	token, _ := os.LookupEnv("AWS_SESSION_TOKEN")

	vars.id = id
	vars.secret = secret
	vars.region = region
	vars.snsTopic = snsTopic
	vars.token = token

	return vars, nil
}

// TTSforrmat formats a string for text-to-speech by removing
// most parantheticals and references.
func TTSformat(s string) string {
	re := regexp.MustCompile(`\n{2,}`)
	s = re.ReplaceAllString(s, "\n\nNext.\n\n")

	// text transformations to improve readability

	// remove parenthetical
	// some are legitimate, but parenthetical in scientific papers
	// generally these are more important, and it is better
	// from a perspective of concentration to remove all
	// parenthetical
	re = regexp.MustCompile(`( ?)\([^)]+ [^)]+\)( ?)`)
	s = re.ReplaceAllString(s, "$1$2")

	// remove number references
	// note \-– for two forms of dashes
	re = regexp.MustCompile(`( ?)\([0-9, \-–]+\)( ?)`)
	s = re.ReplaceAllString(s, "$1$2")

	// remove NCT identifiers without spaces
	re = regexp.MustCompile(`( ?)\(NCT[0-9, ]{5,100}\)( ?)`)
	s = re.ReplaceAllString(s, " ")

	// numbers except for zero, optionally with comma number,
	// even without parentheses
	// that are preceded by lower case letters,
	// are mostly references
	// note \-– for two forms of dashes
	re = regexp.MustCompile(`( ?)[a-z]{2,}[1-9]{1,3}([\-–,][1-9]{1,3})?\.( ?)`)
	s = re.ReplaceAllString(s, "$1$2 ")

	re = regexp.MustCompile(`,[1-9]{1,3}([-–,][1-9]{1,3})`)
	s = re.ReplaceAllString(s, ",")

	re = regexp.MustCompile(`\.[1-9]{1,3}([-–,][1-9]{1,3})`)
	s = re.ReplaceAllString(s, ".")

	// finally, normalize white space
	re = regexp.MustCompile(` {2,}`)
	s = re.ReplaceAllString(s, " ")
	re = regexp.MustCompile(` ,`)
	s = re.ReplaceAllString(s, ",")
	re = regexp.MustCompile(` \.`)
	s = re.ReplaceAllString(s, ".")

	return (s)
}
