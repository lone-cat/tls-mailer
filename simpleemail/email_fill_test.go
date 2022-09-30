package simpleemail

import (
	"os"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	mail := NewEmptyEmail()
	if mail.GetText() != "" {
		t.Errorf(`"Text" is "%s" in empty email`, mail.GetText())
	}
	if mail.GetHtml() != "" {
		t.Errorf(`"Html" is "%s" in empty email`, mail.GetHtml())
	}
	if mail.GetSubject() != "" {
		t.Errorf(`"Subject" is "%s" in empty email`, mail.GetSubject())
	}
	if len(mail.GetFrom()) > 0 {
		t.Errorf(`"From" contains %#v instead of empty list`, mail.GetFrom())
	}
	if len(mail.GetTo()) > 0 {
		t.Errorf(`"To" contains %#v instead of empty list`, mail.GetTo())
	}
	if len(mail.GetCc()) > 0 {
		t.Errorf(`"Cc" contains %#v instead of empty list`, mail.GetCc())
	}
	if len(mail.GetBcc()) > 0 {
		t.Errorf(`"Bcc" contains %#v instead of empty list`, mail.GetBcc())
	}
	if len(mail.(*email).GetEmbedded().ExtractPartsSlice()) > 0 {
		t.Error(`"Embedded" list is not empty in empty email`)
	}
	if len(mail.(*email).GetAttachments().ExtractPartsSlice()) > 0 {
		t.Error(`"Attachments" list is not empty in empty email`)
	}
}

func TestFill(t *testing.T) {
	mail, err := NewEmptyEmail().WithText(text)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.WithHtml(html)
	if err != nil {
		t.Fatal(err)
	}

	mail = mail.WithSubject(subject)

	mail, err = mail.WithFrom(from...)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.WithTo(to...)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.WithCc(cc...)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.WithBcc(bcc...)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		t.Fatal(err)
	}

	mail, err = mail.
		WithAttachedFile(attached)
	if err != nil {
		t.Fatal(err)
	}

	if mail.GetText() != text {
		t.Errorf(`"Text" filled incorrectly. expected %#v, got %#v`, text, mail.GetText())
	}
	if mail.GetHtml() != html {
		t.Errorf(`"Html" filled incorrectly. expected %#v, got %#v`, html, mail.GetHtml())
	}
	if mail.GetSubject() != subject {
		t.Errorf(`"Subject" filled incorrectly. expected %#v, got %#v`, subject, mail.GetSubject())
		t.Fail()
	}
	if !addressSlicesEqual(mail.GetFrom(), from) {
		t.Errorf(`"From" filled incorrectly. expected %#v, got %#v`, from, mail.GetFrom())
	}
	if !addressSlicesEqual(mail.GetTo(), to) {
		t.Errorf(`"To" filled incorrectly. expected %#v, got %#v`, to, mail.GetTo())
	}
	if !addressSlicesEqual(mail.GetCc(), cc) {
		t.Errorf(`"Cc" filled incorrectly. expected %#v, got %#v`, cc, mail.GetCc())
	}
	if !addressSlicesEqual(mail.GetBcc(), bcc) {
		t.Errorf(`"Bcc" filled incorrectly. expected %#v, got %#v`, bcc, mail.GetBcc())
	}
	dataBytes, err := os.ReadFile(embedded)
	if err != nil {
		t.Fatal(err)
	}
	embeddedPartsSlice := mail.(*email).GetEmbedded().ExtractPartsSlice()
	if !reflect.DeepEqual(embeddedPartsSlice[0].GetBody(), dataBytes) {
		t.Errorf(`"Embedded body" filled incorrectly. expected %#v, got %#v`, string(dataBytes), embeddedPartsSlice[0].GetBody())
	}
	dataBytes, err = os.ReadFile(attached)
	if err != nil {
		t.Fatal(err)
	}
	attachmentsPartsSlice := mail.(*email).GetAttachments().ExtractPartsSlice()
	if !reflect.DeepEqual(attachmentsPartsSlice[0].GetBody(), dataBytes) {
		t.Errorf(`"Attached body" filled incorrectly. expected %#v, got %#v`, string(dataBytes), attachmentsPartsSlice[0].GetBody())
	}
}
