package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/urfave/cli"
)

const (
	strlen        = 20
	loginTemplate = `[default]
aws_access_key_id={{.Key}}
aws_secret_access_key={{.Secret}}`
)

type login struct {
	Key    string
	Secret string
}

func duplicationID() string {
	b := make([]byte, strlen)
	rand.Read(b)
	en := base64.StdEncoding // or URLEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	return string(d)
}

func main() {
	// Get user $HOME dir
	usr, err := user.Current()
	if err != nil {
		log.Println("Error getting user home dir", err)
	}

	app := cli.NewApp()
	app.Version = "0.2.1"
	app.Name = "SQS PUSH"
	app.Usage = "push stdin into sqs queue"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cred-path, c",
			Value: usr.HomeDir + "/.aws/credentials",
			Usage: "AWS credentials path",
		},
		cli.StringFlag{
			Name:  "queue, q",
			Usage: "SQS queue URL",
		},
		cli.StringFlag{
			Name:  "region, r",
			Usage: "SQS queue REGION",
		},
		cli.StringFlag{
			Name:  "profile, p",
			Value: "",
			Usage: "AWS login profile",
		},
		cli.Int64Flag{
			Name:  "retries",
			Value: 2,
			Usage: "SQS queue retries",
		},
		cli.StringFlag{
			Name:  "group, g",
			Value: duplicationID(),
			Usage: "SQS queue MessageGroupId",
		},
	}
	app.Action = func(c *cli.Context) error {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			cli.ShowAppHelp(c)
		} else {
			payload := ""
			s := bufio.NewScanner(os.Stdin)
			for s.Scan() {
				payload = payload + s.Text()
			}

			log.Println("Sending: ", payload)

			sess := session.New(&aws.Config{
				Region:      aws.String(c.String("region")),
				Credentials: credentials.NewSharedCredentials(c.String("cred-path"), c.String("profile")),
				MaxRetries:  aws.Int(c.Int("retries")),
			})

			svc := sqs.New(sess)

			// Send message
			sendParams := &sqs.SendMessageInput{
				MessageBody: aws.String(payload),           // Required
				QueueUrl:    aws.String(c.String("queue")), // Required
				// DelaySeconds:           aws.Int64(3),         // (optional) ~ 900s (15 minutes)
				MessageGroupId:         aws.String(c.String("group")),
				MessageDeduplicationId: aws.String(duplicationID()),
			}
			sendResp, err := svc.SendMessage(sendParams)
			if err != nil {
				fmt.Println("err")
				fmt.Println(err)
			}
			fmt.Printf("[Send message] \n%v \n\n", sendResp)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login aws credentials",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "aws access key id",
				},
				cli.StringFlag{
					Name:  "secret",
					Usage: "aws secret access key",
				},
			},
			Action: func(c *cli.Context) error {
				// Create directory if not exists
				if _, err := os.Stat(usr.HomeDir + "/.aws"); os.IsNotExist(err) {
					os.Mkdir(usr.HomeDir+"/.aws", os.ModePerm)
				}

				f, err := os.Create(usr.HomeDir + "/.aws/credentials")
				if err != nil {
					log.Println("create file: ", err)
					return nil
				}

				// Writing template
				t := template.Must(template.New("login").Parse(loginTemplate))
				err = t.Execute(f, login{
					Key:    c.String("key"),
					Secret: c.String("secret"),
				})
				if err != nil {
					log.Println("error writing credentials:", err)
					return nil
				}

				fmt.Println("Login Success")

				return nil
			},
		},
	}

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   cat ./payload.json | {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`

	app.Run(os.Args)
}
