package simpleemail_test

import (
	"github.com/lone-cat/tls-mailer/simpleemail"
	"net/mail"
	"testing"
)

var (
	addr1 = mail.Address{Name: `Иванов Иван`, Address: `first@email.addr`}
	addr2 = mail.Address{Name: `Петров Петр`, Address: `second@email.addr`}
	addr3 = mail.Address{Name: `Сидоров Сидор`, Address: `third@email.addr`}
	addr4 = mail.Address{Name: `Иванов Петр`, Address: `fourth@email.addr`}
	addr5 = mail.Address{Name: `Иванов Сидор`, Address: `fifth@email.addr`}
	addr6 = mail.Address{Name: `Петров Иван`, Address: `sixth@email.addr`}
	addr7 = mail.Address{Name: `Петров Сидор`, Address: `seventh@email.addr`}
	addr8 = mail.Address{Name: `Сидоров Иван`, Address: `eighth@email.addr`}
	addr9 = mail.Address{Name: `Сидоров Петр`, Address: `ninth@email.addr`}
)

var (
	from = []mail.Address{addr1, addr2}
	to   = []mail.Address{addr2, addr3, addr4, addr5}
	cc   = []mail.Address{addr6, addr7, addr8}
	bcc  = []mail.Address{addr9, addr1}

	subject = "Какая то странная длинная тема s angliiskimi simvolami так чтобы \r\nточно поместилась только на несколько строк"

	text = `Какой то не менее длинный текст, чтобы он тоже был на несколько строк, но при этом еще длиннее чем предыдущий` + "\r\n" +
		`Кстати, этот текст еще и будет иметь перевод строки. Tak zhe on soderzhit английские буквы )`
	html = `<h1>Здесь длина текста уже не будет иметь значения</h1>`

	embedded = `aaa`

	attached = `bbb`
)

func TestFill(t *testing.T) {
	email := simpleemail.NewEmptyEmail()
	email.GetText()
}
