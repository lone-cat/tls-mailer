package encoding

const (
	Empty           Type = ``
	SevenBit        Type = `7bit`
	EightBit        Type = `8bit`
	Binary          Type = `binary`
	QuotedPrintable Type = `quoted-printable`
	Base64          Type = `base64`
)

type Type string

func (e Type) String() string {
	return string(e)
}
