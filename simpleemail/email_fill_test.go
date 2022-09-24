package simpleemail

import (
	"os"
	"testing"
)

func TestEmpty(t *testing.T) {
	email := NewEmptyEmail()
	if email.GetText() != "" {
		t.Errorf(`"Text" is "%s" in empty email`, email.GetText())
	}
	if email.GetHtml() != "" {
		t.Errorf(`"Html" is "%s" in empty email`, email.GetHtml())
	}
	if email.GetSubject() != "" {
		t.Errorf(`"Subject" is "%s" in empty email`, email.GetSubject())
	}
	if len(email.GetFrom()) > 0 {
		t.Errorf(`"From" contains %#v instead of empty list`, email.GetFrom())
	}
	if len(email.GetTo()) > 0 {
		t.Errorf(`"To" contains %#v instead of empty list`, email.GetTo())
	}
	if len(email.GetCc()) > 0 {
		t.Errorf(`"Cc" contains %#v instead of empty list`, email.GetCc())
	}
	if len(email.GetBcc()) > 0 {
		t.Errorf(`"Bcc" contains %#v instead of empty list`, email.GetBcc())
	}
	if len(email.GetEmbedded().ExtractPartsSlice()) > 0 {
		t.Error(`"Embedded" list is not empty in empty email`)
	}
	if len(email.GetAttachments().ExtractPartsSlice()) > 0 {
		t.Error(`"Attachments" list is not empty in empty email`)
	}
}

func TestFill(t *testing.T) {
	email := NewEmptyEmail().
		WithText(text).
		WithHtml(html).
		WithSubject(subject).
		WithFrom(from).
		WithTo(to).
		WithCc(cc).
		WithBcc(bcc)
	email, err := email.
		WithEmbeddedFile(cid, embedded)
	if err != nil {
		t.Fatal(err)
	}

	email, err = email.
		WithAttachedFile(attached)
	if err != nil {
		t.Fatal(err)
	}

	if email.GetText() != text {
		t.Errorf(`"Text" filled incorrectly. expected %#v, got %#v`, text, email.GetText())
	}
	if email.GetHtml() != html {
		t.Errorf(`"Html" filled incorrectly. expected %#v, got %#v`, html, email.GetHtml())
	}
	if email.GetSubject() != subject {
		t.Errorf(`"Subject" filled incorrectly. expected %#v, got %#v`, subject, email.GetSubject())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetFrom(), from) {
		t.Errorf(`"From" filled incorrectly. expected %#v, got %#v`, from, email.GetFrom())
	}
	if !addressSlicesEqual(email.GetTo(), to) {
		t.Errorf(`"To" filled incorrectly. expected %#v, got %#v`, to, email.GetTo())
	}
	if !addressSlicesEqual(email.GetCc(), cc) {
		t.Errorf(`"Cc" filled incorrectly. expected %#v, got %#v`, cc, email.GetCc())
	}
	if !addressSlicesEqual(email.GetBcc(), bcc) {
		t.Errorf(`"Bcc" filled incorrectly. expected %#v, got %#v`, bcc, email.GetBcc())
	}
	dataBytes, err := os.ReadFile(embedded)
	if err != nil {
		t.Fatal(err)
	}
	embeddedPartsSlice := email.GetEmbedded().ExtractPartsSlice()
	if embeddedPartsSlice[0].GetBody() != string(dataBytes) {
		t.Errorf(`"Embedded body" filled incorrectly. expected %#v, got %#v`, string(dataBytes), embeddedPartsSlice[0].GetBody())
	}
	dataBytes, err = os.ReadFile(attached)
	if err != nil {
		t.Fatal(err)
	}
	attachmentsPartsSlice := email.GetAttachments().ExtractPartsSlice()
	if attachmentsPartsSlice[0].GetBody() != string(dataBytes) {
		t.Errorf(`"Attached body" filled incorrectly. expected %#v, got %#v`, string(dataBytes), attachmentsPartsSlice[0].GetBody())
	}
}
