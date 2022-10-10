package awslib

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type RegistrationTemplate struct {
	Code string
}

func SendEmail(
	act models.AuthTokenActivation,
	config lib.Configuration,
) (string, error) {

	// get the logo from s3
	awsconf := config.GetAWSConfig()

	sess, err := session.NewSession(&awsconf)
	if err != nil {
		return "", err
	}

	// Create an SES session.
	sessvc := ses.New(sess)

	path, _ := filepath.Abs("awslib/html/Register.html")
	tmpl := template.Must(template.ParseFiles(path))

	var s string
	regg := RegistrationTemplate{
		Code: act.AuthPin,
	}
	buff := bytes.NewBufferString(s)
	err = tmpl.Execute(buff, regg)

	if err != nil {
		return "", err
	}

	// Assemble the email. and send from tavernregister@gmail.com
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(act.AuthEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(buff.String()),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String("Thank you for registering with Tavern!" +
						" Please use this code to complete your registration: " + act.AuthPin),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Tavern Registration"),
			},
		},
		Source: aws.String(" " + config.GetEmailConfig().Username + " "),
	}

	// Attempt to send the email.
	r, err := sessvc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}
	return r.String(), nil

}
