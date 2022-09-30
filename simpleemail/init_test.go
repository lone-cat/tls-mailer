package simpleemail

import (
	"fmt"
	"github.com/lone-cat/stackerrors"
	"net/mail"
)

var cid = `cid`

var (
	addr1 = &mail.Address{Name: `Иванов Иван`, Address: `first@email.addr`}
	addr2 = &mail.Address{Name: `Петров Петр`, Address: `second@email.addr`}
	addr3 = &mail.Address{Name: `Сидоров Сидор`, Address: `third@email.addr`}
	addr4 = &mail.Address{Name: `Иванов Петр`, Address: `fourth@email.addr`}
	addr5 = &mail.Address{Name: `Иванов Сидор`, Address: `fifth@email.addr`}
	addr6 = &mail.Address{Name: `Петров Иван`, Address: `sixth@email.addr`}
	addr7 = &mail.Address{Name: `Петров Сидор`, Address: `seventh@email.addr`}
	addr8 = &mail.Address{Name: `Сидоров Иван`, Address: `eighth@email.addr`}
	addr9 = &mail.Address{Name: `Сидоров Петр`, Address: `ninth@email.addr`}
)

var (
	from = []*mail.Address{addr1, addr2}
	to   = []*mail.Address{addr2, addr3, addr4, addr5}
	cc   = []*mail.Address{addr6, addr7, addr8}
	bcc  = []*mail.Address{addr9, addr1}

	subject = "Какая то странная длинная тема s angliiskimi simvolami так чтобы \r\nточно поместилась только на несколько строк"

	text = `Какой то не менее длинный текст, чтобы он тоже был на несколько строк, но при этом еще длиннее чем предыдущий` + "\r\n" +
		`Кстати, этот текст еще и будет иметь перевод строки. Tak zhe on soderzhit английские буквы )`
	html = `<h1>Здесь длина текста уже не будет иметь значения</h1>`

	embedded = `../test_attachments/image1.jpg`

	attached = `../test_attachments/image2.jpg`
)

var emailsForTest = make([]Email, 0)

func init() {
	a := make(map[string][]string)
	b := make(map[string][]string)
	b = a
	fmt.Println(&a == &b)
	stackerrors.SetDebugMode(true)

	ml := NewEmptyEmail()
	emailsForTest = append(emailsForTest, ml)

	ml, err := ml.WithFrom(from...)
	if err != nil {
		panic(err)
	}

	emailsForTest = append(emailsForTest, ml)

	ml, err = ml.WithTo(to...)
	if err != nil {
		panic(err)
	}

	emailsForTest = append(emailsForTest, ml)

	ml, err = ml.WithCc(cc...)
	if err != nil {
		panic(err)
	}

	emailsForTest = append(emailsForTest, ml)

	ml, err = ml.WithBcc(bcc...)
	if err != nil {
		panic(err)
	}

	emailsForTest = append(emailsForTest, ml)

	ml = ml.WithSubject(subject)
	emailsForTest = append(emailsForTest, ml)

	mail2, err := ml.
		WithText(text)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail2)

	mail2, err = ml.
		WithHtml(html)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail2)

	mail2, err = ml.WithEmbeddedFile(cid, embedded)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail2)

	mail2, err = ml.WithAttachedFile(attached)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail2)

	mail2, err = ml.
		WithText(text)
	if err != nil {
		panic(err)
	}
	mail3, err := mail2.
		WithHtml(html)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	mail3, err = mail2.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	mail3, err = mail2.
		WithAttachedFile(attached)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	mail2, err = ml.
		WithHtml(html)
	if err != nil {
		panic(err)
	}
	mail3, err = mail2.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	mail3, err = mail2.
		WithAttachedFile(attached)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	mail2, err = ml.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		panic(err)
	}
	mail3, err = mail2.
		WithAttachedFile(attached)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

	tmpml, err := ml.
		WithText(text)
	if err != nil {
		panic(err)
	}
	mail3, err = tmpml.
		WithHtml(html)
	if err != nil {
		panic(err)
	}

	mail3, err = mail3.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		panic(err)
	}
	mail3, err = mail3.
		WithAttachedFile(attached)
	if err != nil {
		panic(err)
	}
	emailsForTest = append(emailsForTest, mail3)

}
