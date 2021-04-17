package email

import (
	"bytes"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(name string, email string, token string, tmpl string) error {
	m := gomail.NewMessage()
	m.SetHeader("MIME-version", "1.0")
	m.SetHeader("charset", "UTF-8")
	m.SetHeader("From", m.FormatAddress("noreply@blobber.tk", "ELECT Team"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify and Set Password for your ELECT account.")

	var body bytes.Buffer

	t, err := template.ParseFiles("email/" + tmpl)
	if err != nil {
		return err
	}

	err = t.Execute(&body, map[string]string{
		"name":  name,
		"token": token,
	})
	if err != nil {
		return err
	}

	m.SetBody("text/html", string(body.Bytes()))

	d := gomail.NewDialer("smtp-pulse.com", 587, os.Getenv("SENDPULSE_EMAIL"), os.Getenv("SENDPULSE_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendOTPEmail(email string, otp string, tmpl string) error {
	m := gomail.NewMessage()
	m.SetHeader("MIME-version", "1.0")
	m.SetHeader("charset", "UTF-8")
	m.SetHeader("From", m.FormatAddress("noreply@blobber.tk", "ELECT Team"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "OTP for Login.")

	var body bytes.Buffer

	t, err := template.ParseFiles("email/" + tmpl)
	if err != nil {
		return err
	}

	err = t.Execute(&body, map[string]string{
		"otp": otp,
	})
	if err != nil {
		return err
	}

	m.SetBody("text/html", string(body.Bytes()))

	d := gomail.NewDialer("smtp-pulse.com", 587, os.Getenv("SENDPULSE_EMAIL"), os.Getenv("SENDPULSE_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
