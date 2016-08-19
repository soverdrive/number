package number

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// determine number of float precission, default set to 15 digit
func ProperFloat64(n float64, p ...int) float64 {
	var t float64

	if len(p) > 0 {
		fs := fmt.Sprintf("%s%d%s", "%.", p[0], "f")
		t, _ = strconv.ParseFloat(fmt.Sprintf(fs, n), 64)
	} else {
		t, _ = strconv.ParseFloat(fmt.Sprintf("%.15f", n), 64)
	}

	return t
}

// convert any number to float64 can also take parameter how many digit precission for float64
func NumberToFloat64(n interface{}, p ...int) (float64, error) {
	var t float64

	switch n.(type) {
	case int:
		t = float64(n.(int))
	case int64:
		t = float64(n.(int64))
	case float32:
		t = float64(n.(float32))
	case float64:
		if len(p) > 0 {
			t = ProperFloat64(n.(float64), p[0])
		} else {
			t = ProperFloat64(n.(float64))
		}
	default:
		return t, errors.New("Convert failed, not digit(s)")
	}

	return t, nil
}

// convert any number to int64
func NumberToInt64(n interface{}) (int64, error) {
	var t int64

	switch n.(type) {
	case int:
		t = int64(n.(int))
	case int64:
		t = n.(int64)
	case float32:
		t = int64(n.(float32))
	case float64:
		t = int64(ProperFloat64(n.(float64)))
	default:
		return t, errors.New("Convert failed, not digit(s)")
	}

	return t, nil
}

// seperate string every d-character(s), start with offset (default offset is 0)
func separator(s, r string, d, ds int, buff *bytes.Buffer) {
	var l = len(s)
	var h = ds

	for i := 0; i < l; i++ {
		if h == d {
			buff.WriteString(r)
			h = 0
		}
		h++
		buff.WriteByte(s[i])
	}
}

// seperate string every n-character(s)
func StringSeparator(s, r string, d int) string {
	var b bytes.Buffer
	separator(s, r, d, 0, &b)
	return b.String()
}

// seperate float64 by each d-digits, considering minus sign and setting offset
func float64Separator(t float64, buff *bytes.Buffer, r string, d int) error {
	var isminus bool
	var offset int

	if t < 0 {
		isminus = true
		offset = 1
		buff.WriteByte('-')
	}

	var ts = strconv.FormatFloat(t, 'f', -1, 64)
	var l int
	var isdecimal bool

	l = strings.Index(ts, ".") + 1
	if l > 0 {
		isdecimal = true
	}
	if isminus {
		l--
	}

	if l <= 0 {
		l = len(ts)
	}

	var s_offset int
	s_offset = d - ((l - offset) % d)

	if s_offset == d {
		s_offset = 0
	}

	separator(ts[offset:l], r, d, s_offset, buff)

	if isdecimal {
		if r == "." {
			buff.WriteByte(',')
			l++
		}

		buff.Write([]byte(ts[l:]))
	}

	return nil
}

// private function of money formatting style
func moneyFormat(t float64, buff *bytes.Buffer, r string, d int) error {
	e := float64Separator(t, buff, r, d)
	return e
}

// public function of money formatting style
func MoneyFormat(n interface{}) (string, error) {
	var buff bytes.Buffer

	t, e := NumberToFloat64(n)
	if e != nil {
		return "", e
	}

	e = moneyFormat(t, &buff, ".", 3)
	if e != nil {
		return "", e
	}

	return buff.String(), nil
}

// IDR prefix and suffix styles
const (
	RP_PREFIX       = iota // put prefix "Rp" [default]
	RP_SPACE_PREFIX        // put prefix "Rp "
	RP_DOT_PREFIX          // put prefix "Rp."
	RP_DIGIT_SUFFIX        // put suffix ",-"
	RP_DASH_SUFFIX         // put suffix ",00" [default]
)

// give trailing tz-zero(s) to float64
func float64TrailingZeros(f float64, tz int) float64 {
	fs := fmt.Sprintf("%s%d%s", "%.0", tz, "f")
	t, _ := strconv.ParseFloat(fmt.Sprintf(fs, f), 64)
	return t
}

// give prefix IDR
func idrPrefix(buff *bytes.Buffer, m int) {
	if m == RP_SPACE_PREFIX {
		buff.Write([]byte("Rp "))
	} else if m == RP_DOT_PREFIX {
		buff.Write([]byte("Rp."))
	} else if m == RP_PREFIX {
		buff.Write([]byte("Rp"))
	}
}

// give suffix IDR
func idrSuffix(buff bytes.Buffer, m int) bytes.Buffer {
	var buff_t bytes.Buffer

	if m == RP_DASH_SUFFIX {
		ts := buff.String()
		t := strings.Index(ts, ",") + 1
		buff_t.WriteString(ts[:t])
		buff_t.Write([]byte(",-"))
	} else if m == RP_DIGIT_SUFFIX {
		buff_t = buff
	}

	return buff_t
}

// Rupiah local format
// need to configure how to read param, still give default param
func IDR(n interface{}, p ...int) (string, error) {
	var buff bytes.Buffer
	var prefix int
	var suffix int

	prefix = getParam(p)

	idrPrefix(&buff, RP_PREFIX)

	t, e := NumberToFloat64(n, 2)
	if e != nil {
		return "", e
	}

	t = float64TrailingZeros(t, 2)

	e = moneyFormat(t, &buff, ".", 3)
	if e != nil {
		return "", e
	}

	buff = idrSuffix(buff, RP_DIGIT_SUFFIX)

	return buff.String(), nil
}
