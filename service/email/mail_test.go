package email

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

func TestSendEmail(t *testing.T) {
	err := SendEmail("邮箱验证码", fmt.Sprintf(VerificationCode, "666666"), []string{"710816334@qq.com"})
	t.Log(err)
}

func TestTmpl(t *testing.T) {
	params := map[string]interface{}{
		"ID":   1,
		"Name": "zxz",
	}
	tmpl, err := template.New("test").Parse("The name for student {{.ID}} is {{.Name}}")
	if err != nil {
		t.Fatal(err.Error())
	}
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, params)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(tmplBytes.String())
}
