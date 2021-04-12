package email

import (
	"github.com/go-gomail/gomail"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/music-recommend/config"
)

const (
	fromName    = "From"
	sendTo      = "To"
	sendSubject = "subject"
	sendBody    = "text/html"

	ServerMsg = "STA音乐推荐系统消息 Level[%s]"
)

type emailService struct {
	helper *gomail.Dialer
}

var PubEmailService emailService

func InitEmailService() {
	PubEmailService = emailService{}
	cfg := config.GlobalConfig().EmailConfig
	PubEmailService.helper = gomail.NewDialer(cfg.Host, cfg.Port, cfg.Email, cfg.Pwd)
}

func (es *emailService) newMessage(subject, body string, mailTo ...string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader(fromName, config.GlobalConfig().EmailConfig.Email, config.GlobalConfig().ServerName)
	message.SetHeader(sendTo, mailTo...)
	message.SetHeader(sendSubject, subject)
	message.SetBody(sendBody, body)
	return message
}

func (es *emailService) SendEmail(subject, body string, mailTo ...string) error {
	err := es.helper.DialAndSend(es.newMessage(subject, body, mailTo...))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
