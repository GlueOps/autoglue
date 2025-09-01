package smtp

import (
	"bytes"
	"time"

	"github.com/glueops/autoglue/internal/assets"
	"github.com/glueops/autoglue/internal/funcs"
	"github.com/wneessen/go-mail"

	htmlTemplate "html/template"
	textTemplate "text/template"
)

const defaultTimeout = 10 * time.Second

type Mailer struct {
	client *mail.Client
	from   string
}

func NewMailer(host string, port int, username, password, from string) (*Mailer, error) {
	opts := []mail.Option{
		mail.WithTimeout(defaultTimeout),
		mail.WithPort(port),
		// IMPORTANT for Mailpit/local dev: no TLS
		mail.WithTLSPolicy(mail.NoTLS),
	}

	if username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthLogin),
			mail.WithUsername(username),
			mail.WithPassword(password),
		)
	}

	client, err := mail.NewClient(host, opts...)
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		client: client,
		from:   from,
	}

	return mailer, nil
}

func (m *Mailer) Send(recipient string, data any, patterns ...string) error {
	for i := range patterns {
		patterns[i] = "emails/" + patterns[i]
	}
	msg := mail.NewMsg()

	err := msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(m.from)
	if err != nil {
		return err
	}

	ts, err := textTemplate.New("").Funcs(funcs.TemplateFuncs).ParseFS(assets.EmbeddedFiles, patterns...)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = ts.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())

	plainBody := new(bytes.Buffer)
	err = ts.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())

	if ts.Lookup("htmlBody") != nil {
		ts, err := htmlTemplate.New("").Funcs(funcs.TemplateFuncs).ParseFS(assets.EmbeddedFiles, patterns...)
		if err != nil {
			return err
		}

		htmlBody := new(bytes.Buffer)
		err = ts.ExecuteTemplate(htmlBody, "htmlBody", data)
		if err != nil {
			return err
		}

		msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())
	}

	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)

		if nil == err {
			return nil
		}

		if i != 3 {
			time.Sleep(2 * time.Second)
		}
	}

	return err
}
