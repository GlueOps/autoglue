package authn

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
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

func mustInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func adminCount(except *uuid.UUID) (int64, error) {
	q := db.DB.Model(&models.User{}).Where(`role = ?`, "admin")
	if except != nil {
		q = q.Where("id <> ?", *except)
	}
	var n int64
	err := q.Count(&n).Error
	return n, err
}

func requireGlobalAdmin(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	ctx := middleware.GetAuthContext(r)
	if ctx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	var me models.User
	if err := db.DB.Select("id, role").First(&me, "id = ?", ctx.UserID).Error; err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	if me.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil, false
	}
	return &me, true
}

func asUserOut(u models.User) userOut {
	return userOut{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Role:          string(u.Role),
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
