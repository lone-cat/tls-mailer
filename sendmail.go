package tls_mailer

/*
 * from the code sight main difference from "net/smtp"."SendMail()" is mostly in "Dial" function. Other parts just a bit
 * fixed to use unexported field "net/smtp"."Client"."ext" and unexported function "net/smtp"."Client"."hello()".
 * Unexported function has no good way to be used (if you know better way - please contact me), so I used exported
 * "net/smtp"."Client"."Hello()" function instead. It is not exact replacement, but I guess difference is acceptable. At
 * least more acceptable than excluding hello call at all.
 *
 * Function-helper "tlsDial" is replacement for original "net/smtp"."Dial()" function with same signature, but it calls
 * not "net"."Dial()" as "net/smtp"."Dial()" does, but "tls"."Dial()".
 *
 * Function-helper "validateLine()" is direct copy from "net/smtp" package.
 *
 * Function-helper "getExtFromClient()" is made to extract "ext" property from "net/smtp"."Client" struct using
 * reflection. This "ext" property is used in native "net/smtp" package and can't be somehow easily got from struct
 * another way.
 *
 * Some comments were added to "SendMail" function to explain changes that were made
 */

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"reflect"
	"strings"
)

const (
	extProperty = `ext`
)

var localName = `localhost`

func SetLocalName(localHostName string) error {
	err := validateLine(localName)
	if err != nil {
		return err
	}
	localName = localHostName
	return nil
}

func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {

	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	c, err := tlsDial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	// Hello is not needed if local name is localhost
	err = c.Hello(localName)
	if err != nil {
		return err
	}

	// starttls doesn't need to be used if full connection is encrypted.
	/*
		if ok, _ := c.Extension("STARTTLS"); ok {
			config := &tls.Config{ServerName: host}
			if err = c.StartTLS(config); err != nil {
				return err
			}
		}
	*/

	// this block was added because smtp.Client.ext property can't be directly accessed in next block
	ext, err := getExtFromClient(c)
	if err != nil {
		return err
	}

	if a != nil && ext != nil {
		if _, ok := ext["AUTH"]; !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func tlsDial(addr string) (client *smtp.Client, err error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		//InsecureSkipVerify: false, // unnecessary default here, just indicates that tls connection is verified
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return
	}

	client, err = smtp.NewClient(conn, host)
	return
}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func getExtFromClient(client *smtp.Client) (map[string]string, error) {
	if client == nil {
		return nil, errors.New(`nil client passed`)
	}

	extRefl := reflect.ValueOf(client).Elem().FieldByName(extProperty)

	if extRefl == (reflect.Value{}) {
		return nil, errors.New(`"` + extProperty + `" property does not exist`)
	}

	if extRefl.IsNil() {
		return nil, nil
	}

	if extRefl.Kind() != reflect.Map {
		return nil, errors.New(`"` + extProperty + `" is not map`)
	}

	result := make(map[string]string)
	for _, keyValue := range extRefl.MapKeys() {
		if keyValue.Kind() != reflect.String {
			return nil, errors.New(`not string map key`)
		}
		value := extRefl.MapIndex(keyValue)
		if value.Kind() != reflect.String {
			return nil, errors.New(`not string map value`)
		}
		result[keyValue.String()] = value.String()
	}

	return result, nil
}
