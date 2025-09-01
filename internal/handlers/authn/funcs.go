package authn

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	appsmtp "github.com/glueops/autoglue/internal/smtp"
	"github.com/google/uuid"
)

func randomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func issuePasswordReset(userID uuid.UUID, email string) (string, error) {
	tok, err := randomToken(32)
	if err != nil {
		return "", err
	}
	pr := models.PasswordReset{
		UserID:    userID,
		Token:     tok, // consider storing hash in prod
		ExpiresAt: time.Now().Add(resetTTL),
		Used:      false,
	}
	if err := db.DB.Create(&pr).Error; err != nil {
		return "", err
	}
	return tok, nil
}

func issueEmailVerification(userID uuid.UUID, email string) (string, error) {
	tok, err := randomToken(32)
	if err != nil {
		return "", err
	}
	ev := models.EmailVerification{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tok, // consider storing hash in prod
		ExpiresAt: time.Now().Add(verifyTTL),
		Used:      false,
	}
	if err := db.DB.Create(&ev).Error; err != nil {
		return "", err
	}
	return tok, nil
}

func sendEmail(to, subject, body string) error {
	// integrate with your provider here
	fmt.Printf("Sending email to: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Content-Type: text/html; charset=UTF-8\n")
	fmt.Printf("%s\n", body)
	return nil
}

func getMailer() (*appsmtp.Mailer, error) {
	mailerOnce.Do(func() {
		if !config.SMTPEnabled() {
			mailerErr = fmt.Errorf("smtp disabled")
			return
		}
		mailer, mailerErr = appsmtp.NewMailer(
			config.SMTPHost(),
			config.SMTPPort(),
			config.SMTPUsername(),
			config.SMTPPassword(),
			config.SMTPFrom(),
		)
	})
	return mailer, mailerErr
}

func sendTemplatedEmail(to string, templateFile string, data any) error {
	m, err := getMailer()
	if err != nil {
		// fail soft if smtp is disabled; return nil so API UX isn't blocked
		return nil
	}
	return m.Send(to, data, templateFile)
}
