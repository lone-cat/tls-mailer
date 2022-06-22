package simpleemail_test

var exampleShort = "" +
	"Content-Transfer-Encoding: quoted-printable\r\n" +
	"Content-Type: text/plain; charset=utf-8\r\n" +
	"From: =?utf-8?q?=D0=9E=D1=82_=D0=9A=D0=BE=D0=B3=D0=BE-=D0=A2=D0=BE?= <some@address.com>\r\n" +
	"GetSubject: =?utf-8?q?=D0=9A=D0=B0=D0=BA=D0=B0=D1=8F_=D1=82=D0=BE_=D1=85=D0=BE=D1=80?=\r\n" +
	" =?utf-8?q?=D0=BE=D1=88=D0=B0=D1=8F_=D1=82=D0=B5=D0=BC=D0=B0?=\r\n" +
	"To: =?utf-8?q?=D0=9A=D0=BE=D0=BC=D1=83-=D1=82=D0=BE?= <some-another@address.com>\r\n" +
	"\r\n" +
	"=D0=A5=D0=BE=D1=80=D0=BE=D1=88=D0=B5=D0=B5 =D0=BF=D0=B8=D1=81=D1=8C=D0=BC=\r\n" +
	"=D0=BE\r\n"

var exampleLongText = "" +
	"Content-Transfer-Encoding: quoted-printable\r\n" +
	"Content-Type: text/plain; charset=utf-8\r\n" +
	"Date: Mon, 20 Jun 2022 22:06:56 +0300\r\n" +
	"From: test_portal@melentev.net\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Message-ID: <88b8e287a9eba276950ef6984e8f0a62@melentev.net>\r\n" +
	"GetSubject: =?utf-8?q?=D0=A5=D0=BE=D1=80=D0=BE=D1=88=D0=B0=D1=8F_?=\r\n" +
	" =?utf-8?q?=D1=82=D0=B5=D0=BC=D0=B0_=D0=BD=D0=BE_=D0=BE=D0=BD=D0=B0_?=\r\n" +
	" =?utf-8?q?=D0=BE=D1=87=D0=B5=D0=BD=D1=8C_=D0=B4=D0=BB?=\r\n" +
	" =?utf-8?q?=D0=B8=D0=BD=D0=BD=D0=B0=D1=8F_=D0=B8_=D1=81=D0=BE=D0=B4=D0=B5?=\r\n" +
	" =?utf-8?q?=D1=80=D0=B6=D0=B8=D1=82_=D0=B0=D0=BD=D0=B3?=\r\n" +
	" =?utf-8?q?=D0=BB=D0=B8=D0=B9=D1=81=D0=BA=D0=B8=D0=B5?= symbols!\r\n" +
	"To: Alex <alex@melentev.net>\r\n" +
	"\r\n" +
	"body and image =D0=A5=D0=BE=D1=80=D0=BE=D1=88=D0=B0=D1=8F =D1=82=D0=B5=\r\n" +
	"=D0=BC=D0=B0 =D0=BD=D0=BE =D0=BE=D0=BD=D0=B0 =D0=BE=D1=87=D0=B5=D0=BD=D1=\r\n" +
	"=8C =D0=B4=D0=BB=D0=B8=D0=BD=D0=BD=D0=B0=D1=8F =D0=B8 =D1=81=D0=BE=D0=B4=\r\n" +
	"=D0=B5=D1=80=D0=B6=D0=B8=D1=82 =D0=B0=D0=BD=D0=B3=D0=BB=D0=B8=D0=B9=D1=\r\n" +
	"=81=D0=BA=D0=B8=D0=B5 symbols! =D0=A5=D0=BE=D1=80=D0=BE=D1=88=D0=B0=D1=\r\n" +
	"=8F =D1=82=D0=B5=D0=BC=D0=B0 =D0=BD=D0=BE =D0=BE=D0=BD=D0=B0 =D0=BE=D1=\r\n" +
	"=87=D0=B5=D0=BD=D1=8C =D0=B4=D0=BB=D0=B8=D0=BD=D0=BD=D0=B0=D1=8F =D0=B8 =\r\n" +
	"=D1=81=D0=BE=D0=B4=D0=B5=D1=80=D0=B6=D0=B8=D1=82 =D0=B0=D0=BD=D0=B3=D0=\r\n" +
	"=BB=D0=B8=D0=B9=D1=81=D0=BA=D0=B8=D0=B5 symbols! =D0=A5=D0=BE=D1=80=D0=\r\n" +
	"=BE=D1=88=D0=B0=D1=8F =D1=82=D0=B5=D0=BC=D0=B0 =D0=BD=D0=BE =D0=BE=D0=BD=\r\n" +
	"=D0=B0 =D0=BE=D1=87=D0=B5=D0=BD=D1=8C =D0=B4=D0=BB=D0=B8=D0=BD=D0=BD=D0=\r\n" +
	"=B0=D1=8F =D0=B8 =D1=81=D0=BE=D0=B4=D0=B5=D1=80=D0=B6=D0=B8=D1=82 =D0=B0=\r\n" +
	"=D0=BD=D0=B3=D0=BB=D0=B8=D0=B9=D1=81=D0=BA=D0=B8=D0=B5 symbols! =D0=A5=\r\n" +
	"=D0=BE=D1=80=D0=BE=D1=88=D0=B0=D1=8F =D1=82=D0=B5=D0=BC=D0=B0 =D0=BD=D0=\r\n" +
	"=BE =D0=BE=D0=BD=D0=B0 =D0=BE=D1=87=D0=B5=D0=BD=D1=8C =D0=B4=D0=BB=D0=B8=\r\n" +
	"=D0=BD=D0=BD=D0=B0=D1=8F =D0=B8 =D1=81=D0=BE=D0=B4=D0=B5=D1=80=D0=B6=D0=\r\n" +
	"=B8=D1=82 =D0=B0=D0=BD=D0=B3=D0=BB=D0=B8=D0=B9=D1=81=D0=BA=D0=B8=D0=B5 symb=\r\n" +
	"ols!=D0=A5=D0=BE=D1=80=D0=BE=D1=88=D0=B0=D1=8F =D1=82=D0=B5=D0=BC=D0=B0 =\r\n" +
	"=D0=BD=D0=BE =D0=BE=D0=BD=D0=B0 =D0=BE=D1=87=D0=B5=D0=BD=D1=8C =D0=B4=D0=\r\n" +
	"=BB=D0=B8=D0=BD=D0=BD=D0=B0=D1=8F =D0=B8 =D1=81=D0=BE=D0=B4=D0=B5=D1=80=\r\n" +
	"=D0=B6=D0=B8=D1=82 =D0=B0=D0=BD=D0=B3=D0=BB=D0=B8=D0=B9=D1=81=D0=BA=D0=\r\n" +
	"=B8=D0=B5 symbols!\r\n"

var exampleTextAndHtml = "" +
	"Content-Type: multipart/alternative; boundary=54CF9i9Q\r\n" +
	"Date: Tue, 21 Jun 2022 09:36:56 +0300\r\n" +
	"From: test_portal@melentev.net\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Message-ID: <d99007196434f6731ae382e5b7159525@melentev.net>\r\n" +
	"GetSubject: =?utf-8?Q?=D0=9A=D0=B0=D0=BA=D0=B0=D1=8F-=D1=82=D0=BE?=\r\n" +
	" =?utf-8?Q?_=D1=82=D0=B5=D0=BC=D0=B0?=\r\n" +
	"To: Alex <alex@melentev.net>\r\n" +
	"\r\n" +
	"--54CF9i9Q\r\n" +
	"Content-Type: text/plain; charset=utf-8\r\n" +
	"Content-Transfer-Encoding: quoted-printable\r\n" +
	"\r\n" +
	"some text\r\n" +
	"--54CF9i9Q\r\n" +
	"Content-Type: text/html; charset=utf-8\r\n" +
	"Content-Transfer-Encoding: quoted-printable\r\n" +
	"\r\n" +
	"<h1>html here</h1>\r\n" +
	"--54CF9i9Q--"
