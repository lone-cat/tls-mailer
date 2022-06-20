package simpleemail

import "io"

type splitter struct {
	separator []byte
	len       uint
	writer    io.Writer
	written   uint
}

func NewSplitter(writer io.Writer, separator []byte, len uint) *splitter {
	return &splitter{
		separator: separator,
		len:       len,
		writer:    writer,
	}
}

func (s *splitter) Write(data []byte) (written int, err error) {
	buffer := make([]byte, s.len)
	ln := 0
	var toRead uint
	for len(data) > 0 {
		toRead = s.len - s.written
		if len(data) >= int(toRead) {
			copy(buffer, data[:toRead])
			ln, err = s.writer.Write(buffer[:toRead])
			s.written += uint(ln)
			s.written %= s.len
			if err != nil {
				return
			}
			ln, err = s.writer.Write(s.separator)
			if err != nil {
				return
			}
			data = data[toRead:]
		} else {
			ln = copy(buffer, data)
			ln, err = s.writer.Write(buffer[:ln])
			s.written += uint(ln)
			s.written %= s.len
			if err != nil {
				return
			}
			data = nil
		}
	}

	return
}
