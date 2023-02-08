package mailer

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	htemplate "html/template"
	"path/filepath"
	ttemplate "text/template"

	"github.com/csmith/envflag"
	"github.com/mailgun/mailgun-go/v4"
)

var (
	mailDomain  = flag.String("mailgun-domain", "", "Domain to use when sending mail from mailgun")
	mailSender  = flag.String("mailgun-sender", "", "From e-mail address to use for e-mails")
	apiKey      = flag.String("mailgun-api-key", "", "API key to use for mailgun")
	apiBase     = flag.String("mailgun-api-base", "https://api.eu.mailgun.net/v3", "Base URL for the mailgun API")
	templateDir = flag.String("template-dir", "templates", "Directory for Puzzad e-mail templates")
)

func init() {
	envflag.Parse()
}

func sendConfirmationEmail(ctx context.Context, email, link string) error {
	return sendMail(ctx, email, "Puzzad: Confirm sign up", "signup", map[string]string{
		"Link": link,
	})
}

func sendMail(ctx context.Context, address, subject, template string, data any) error {
	return sendMailWithReplyTo(ctx, address, subject, template, "", data)
}

func sendMailWithReplyTo(ctx context.Context, address, subject, template, replyTo string, data any) error {
	mg := mailgun.NewMailgun(*mailDomain, *apiKey)
	mg.SetAPIBase(*apiBase)

	tt, err := ttemplate.ParseFiles(filepath.Join(*templateDir, fmt.Sprintf("%s.txt.gotpl", template)))
	if err != nil {
		return err
	}

	ht, err := htemplate.ParseFiles(filepath.Join(*templateDir, fmt.Sprintf("%s.html.gotpl", template)))
	if err != nil {
		return err
	}

	twr := &bytes.Buffer{}
	if err = tt.Execute(twr, data); err != nil {
		return err
	}

	hwr := &bytes.Buffer{}
	if err = ht.Execute(hwr, data); err != nil {
		return err
	}

	message := mg.NewMessage(*mailSender, subject, twr.String(), address)
	if replyTo != "" {
		message.SetReplyTo(replyTo)
	}
	message.SetHtml(hwr.String())
	_, _, err = mg.Send(ctx, message)
	return err
}
