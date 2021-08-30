package email

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码
	SPassword string
	// SMTP 服务器地址
	SMTPAddr string
	// SMTP端口
	SMTPPort int
}

func SendEmail(title, body string, emails []string) error {
	var mailConf MailboxConf
	mailConf.Title = title
	//mailConf.Body = `<html> <body> <h3> 您好： </h3> 非常感谢您使用氚平台，您的邮箱验证码为：<br/> <b>` + code + `</b><br/> 此验证码有效期30分钟，请妥善保存。<br/> 如果这不是您本人的操作，请忽略本邮件。<br/> </body> </html> `
	mailConf.Body = body
	mailConf.RecipientList = emails
	mailConf.Sender = `service@dataqin.com`
	mailConf.SPassword = "19WQ2PFqtGv2q245"
	mailConf.SMTPAddr = `smtp.feishu.cn`
	mailConf.SMTPPort = 465

	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender)
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, mailConf.Body)
	//m.Attach()   //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		logrus.Errorf("Send Email Fail, %s", err.Error())
		return err
	}
	return nil
}
