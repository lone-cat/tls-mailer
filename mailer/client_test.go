package mailer_test

import (
	"encoding/json"
	"fmt"
	"github.com/lone-cat/tls-mailer/mailer"
	"net/mail"
	"os"
	"testing"
)

type config struct {
	Server   string `json:"server"`
	Port     uint16 `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

func (c *config) String() string {
	bts, err := json.MarshalIndent(c, ``, `  `)
	if err != nil {
		return err.Error()
	}

	return string(bts)
}

func TestClient(t *testing.T) {
	configBytes, err := os.ReadFile(`./config.json`)
	if err != nil {
		t.Fatal(err)
	}

	var conf *config = &config{}
	err = json.Unmarshal(configBytes, conf)
	if err != nil {
		t.Fatal(err)
	}

	sender, err := mail.ParseAddress(conf.Sender)
	if err != nil {
		t.Fatal(err)
	}

	sender.Name = `Test`

	cl, err := mailer.NewClient(mailer.StartTLSClient, conf.Server, conf.Port, conf.Login, conf.Password, sender)
	if err != nil {
		t.Fatal(err)
	}

	email := cl.Email()

	receiver, err := mail.ParseAddress(conf.Receiver)
	if err != nil {
		t.Fatal(err)
	}

	email, err = email.WithTo(receiver)
	if err != nil {
		t.Fatal(err)
	}

	email, err = email.WithHtml(`<html>some text<img src="cid:img1" /></html>`)
	if err != nil {
		t.Fatal(err)
	}

	email, err = email.WithEmbeddedFile(`img1`, `../test_attachments/image1.jpg`)
	if err != nil {
		t.Fatal(err)
	}

	email, err = email.WithAttachedFile(`../test_attachments/new.txt`)
	if err != nil {
		t.Fatal(err)
	}

	email = email.WithAttachedBytes([]byte(`dfczxz`))

	email, err = email.WithText(`html>some text<img src="cid:img1" /></html>`)
	if err != nil {
		t.Fatal(err)
	}

	email = email.WithoutAttachments()

	dumpEmail, err := json.MarshalIndent(email.Dump(), ``, `  `)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dumpEmail))

	//err = cl.Send(email)
	if err != nil {
		t.Fatal(err)
	}

}
