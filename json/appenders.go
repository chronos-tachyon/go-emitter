package json

import (
	"encoding/base64"
	"fmt"
)

type StringAppender struct {
	Value      string
	EscapeHTML bool
}

func (a StringAppender) String() string {
	return string(a.Append(nil))
}

func (a StringAppender) Append(out []byte) []byte {
	out = append(out, '"')
	for _, ch := range a.Value {
		if esc, found := stringEscapes[ch]; found {
			out = append(out, esc...)
			continue
		}

		if a.EscapeHTML {
			if esc, found := stringEscapesHTML[ch]; found {
				out = append(out, esc...)
				continue
			}
		}

		if ch >= 0xd800 && ch < 0xe000 {
			out = fmt.Appendf(out, "\\u%04x", ch)
			continue
		}

		str := string(ch)
		out = append(out, str...)
	}
	return append(out, '"')
}

var (
	_ fmt.Stringer = StringAppender{}
	_ Appender     = StringAppender{}
)

type BytesAppender struct {
	Value []byte
}

func (a BytesAppender) String() string {
	return string(a.Append(nil))
}

func (a BytesAppender) Append(out []byte) []byte {
	out = append(out, '"')
	n := uint(len(a.Value))
	for i := uint(0); i < n; i += 3 {
		j := min(i+3, n)
		var tmp [4]byte
		base64.StdEncoding.Encode(tmp[:], a.Value[i:j])
		out = append(out, tmp[:]...)
	}
	return append(out, '"')
}

var (
	_ fmt.Stringer = BytesAppender{}
	_ Appender     = BytesAppender{}
)

var stringEscapes = map[rune]string{
	'"':    `\"`,
	'\\':   `\\`,
	'\x00': `\u0000`, // U+0000 NUL
	'\x01': `\u0001`, // U+0001 SOH
	'\x02': `\u0002`, // U+0002 STX
	'\x03': `\u0003`, // U+0003 ETX
	'\x04': `\u0004`, // U+0004 EOT
	'\x05': `\u0005`, // U+0005 ENQ
	'\x06': `\u0006`, // U+0006 ACK
	'\x07': `\u0007`, // U+0007 BEL
	'\b':   `\b`,     // U+0008 BS
	'\t':   `\t`,     // U+0009 HT
	'\n':   `\n`,     // U+000a LF
	'\x0b': `\u000b`, // U+000b VT
	'\f':   `\f`,     // U+000c FF
	'\r':   `\r`,     // U+000d CR
	'\x0e': `\u000e`, // U+000e SO
	'\x0f': `\u000f`, // U+000f SI
	'\x10': `\u0010`, // U+0010 DLE
	'\x11': `\u0011`, // U+0011 DC1
	'\x12': `\u0012`, // U+0012 DC2
	'\x13': `\u0013`, // U+0013 DC3
	'\x14': `\u0014`, // U+0014 DC4
	'\x15': `\u0015`, // U+0015 NAK
	'\x16': `\u0016`, // U+0016 SYN
	'\x17': `\u0017`, // U+0017 ETB
	'\x18': `\u0018`, // U+0018 CAN
	'\x19': `\u0019`, // U+0019 EM
	'\x1a': `\u001a`, // U+001a SUB
	'\x1b': `\u001b`, // U+001b ESC
	'\x1c': `\u001c`, // U+001c FS
	'\x1d': `\u001d`, // U+001d GS
	'\x1e': `\u001e`, // U+001e RS
	'\x1f': `\u001f`, // U+001f US
	'\x7f': `\u007f`, // U+007f DEL
	'\x80': `\u0080`, // U+0080 PAD
	'\x81': `\u0081`, // U+0081 HOP
	'\x82': `\u0082`, // U+0082 BPH
	'\x83': `\u0083`, // U+0083 NBH
	'\x84': `\u0084`, // U+0084 IND
	'\x85': `\u0085`, // U+0085 NEL
	'\x86': `\u0086`, // U+0086 SSA
	'\x87': `\u0087`, // U+0087 ESA
	'\x88': `\u0088`, // U+0088 HTS
	'\x89': `\u0089`, // U+0089 HTJ
	'\x8a': `\u008a`, // U+008a VTS
	'\x8b': `\u008b`, // U+008b PLD
	'\x8c': `\u008c`, // U+008c PLU
	'\x8d': `\u008d`, // U+008d RI
	'\x8e': `\u008e`, // U+008e SS2
	'\x8f': `\u008f`, // U+008f SS3
	'\x90': `\u0090`, // U+0090 DCS
	'\x91': `\u0091`, // U+0091 PU1
	'\x92': `\u0092`, // U+0092 PU2
	'\x93': `\u0093`, // U+0093 STS
	'\x94': `\u0094`, // U+0094 CCH
	'\x95': `\u0095`, // U+0095 MW
	'\x96': `\u0096`, // U+0096 SPA
	'\x97': `\u0097`, // U+0097 EPA
	'\x98': `\u0098`, // U+0098 SOS
	'\x99': `\u0099`, // U+0099 SGCI
	'\x9a': `\u009a`, // U+009a SCI
	'\x9b': `\u009b`, // U+009b CSI
	'\x9c': `\u009c`, // U+009c ST
	'\x9d': `\u009d`, // U+009d OSC
	'\x9e': `\u009e`, // U+009e PM
	'\x9f': `\u009f`, // U+009f APC
}

var stringEscapesHTML = map[rune]string{
	'&': `\u0026`,
	'<': `\u003c`,
	'>': `\u003e`,
}
