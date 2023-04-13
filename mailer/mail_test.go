package mailer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	smtpAuthAddress   = "localhost"
	smtpServerAddress = "localhost:1025"
	mailSenderName    = "Test Send Email"
	mailSenderAddress = "simplify@admin.com"
)

func TestSendEmail(t *testing.T) {
	sender := NewMailSender(
		smtpAuthAddress,
		smtpServerAddress,
		mailSenderName,
		mailSenderAddress,
		"",
	)

	subject := "Test Send Email"
	content := `
		<h3>您好，欢迎注册 ShortNet</h3>
		<p>如果是你本人注册 ShortNet，请点击下面激活账户，若不是请忽略该邮件。</p>
		<p><a href="https://github.com/lightsaid/short-net">激活账户</a></p>
	`

	to := []string{"7zeroc@gmail.com", "abc@163.com"}
	attchFiles := []string{"../README.md"}

	err := sender.SendEmail(subject, content, to, nil, nil, attchFiles)
	require.NoError(t, err)
}
