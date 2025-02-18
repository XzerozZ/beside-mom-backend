package utils

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"bytes"
	"errors"
	"html/template"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendMail(templatePath string, user entities.User, otp string, config configs.Mail) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct {
		Username string
		OTP      string
	}{
		Username: user.Firstname,
		OTP:      otp,
	})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Recovery Your Password")
	m.SetBody("text/html", body.String())
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Host, port, config.Sender, config.Key)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendPasswordMail(templatePath string, user entities.User, password string, config configs.Mail) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct {
		Username string
		Password string
	}{
		Username: user.Firstname,
		Password: password,
	})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Your Password for Beside Mom")
	m.SetBody("text/html", body.String())
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Host, port, config.Sender, config.Key)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(email)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("invalid Email")
	}

	localPart, domain := parts[0], parts[1]
	localPart = strings.ReplaceAll(localPart, ".", "")
	email = localPart + "@" + domain
	return email, nil
}
