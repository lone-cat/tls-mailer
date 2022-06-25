package simpleemail_test

import (
	"fmt"
	"testing"
)

type testStruct struct {
	a string
	b []int
}

func (t testStruct) WithA(a string) testStruct {
	t.a = a
	return t
}

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
	a := testStruct{
		a: `wtf`,
		b: []int{1},
	}

	b := a.WithA(`wtf2`)
	b.b = append(b.b, 2)

	fmt.Printf("%#v\r\n", a)
	fmt.Printf("%#v\r\n", b)
	return
	/*email := simpleemail.NewEmptyEmail().WithTo(to) //.WithAttachedFile(`../test_attachments/image1.jpg`)
	//email, _ = email.WithEmbeddedFile(`some cid`, `../test_attachments/image2.jpg`)
	email, err := simpleemail.Import(email.String())
	fmt.Printf("%#v\r\n", email)
	return
	email, err = simpleemail.Import(exampleTextAndHtml)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(email.GetSubject())
	return

	fmt.Println(exampleTextAndHtml)
	fmt.Println(email.String() == exampleTextAndHtml)*/
}
