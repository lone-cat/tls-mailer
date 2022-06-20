package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"testing"
)

func TestMainFunc(t *testing.T) {
	/*email := simpleemail.NewEmptyEmail()
	//email = email.WithText(`some text`)

	//email = email.WithHtml(`<h1>some text</h1>`)
	email = email.WithSubject(`Моя тема`)

	email = email.WithFrom([]mail.Address{{`Сашка`, `b@a`}})
	email = email.WithTo([]mail.Address{{`Машка`, `b@a`}})

	email, err := email.WithAttachedFile(`../config.json`)
	if err != nil {
		panic(err)
	}
	email, err = email.WithEmbeddedFile(`id`, `../config.json`)
	if err != nil {
		panic(err)
	}
	email, err = email.WithEmbeddedFile(`id`, `../config.json`)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%v\r\n", email.String())
	fmt.Println(email)
	fmt.Println(err)
	*/
	email, err := simpleemail.Import(exampleShort)
	if err != nil {
		panic(err)
	}

	fmt.Println(email.String())
	fmt.Println(exampleShort)
	fmt.Println(email.String() == exampleShort)
}