package number

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ProperFloat64(n float64) float64 {
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.15f", n), 64)
	return t
}

func NumberToFloat64(n interface{}) (float64, bool, error) {
	var f bool
	var t float64

	switch n.(type) {
	case int:
		t = float64(n.(int))
		f = false
	case int64:
		t = float64(n.(int64))
		f = false
	case float32:
		t = float64(n.(float32))
		f = true
	case float64:
		t = ProperFloat64(n.(float64))
		f = true
	default:
		return t, f, errors.New("Convert failed, not digit(s)")
	}

	return t, f, nil
}

func NumberToInt64(n interface{}) (int64, bool, error) {
	var f bool
	var t int64

	switch n.(type) {
	case int:
		t = int64(n.(int))
		f = false
	case int64:
		t = n.(int64)
		f = false
	case float32:
		t = int64(n.(float32))
		f = true
	case float64:
		t = int64(ProperFloat64(n.(float64)))
		f = true
	default:
		return t, f, errors.New("Convert failed, not digit(s)")
	}

	return t, f, nil
}

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

func StringSeparator(s, r string, d int) string {
	var b bytes.Buffer
	separator(s, r, d, 0, &b)
	return b.String()
}

func numberSeparator(n interface{}, buff *bytes.Buffer, r string, d int) error {
	var isfloat bool

	t, isfloat, e := NumberToFloat64(n)
	if e != nil {
		return e
	}

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

	if isfloat {
		l = strings.Index(ts, ".") + 1
		if l > 0 {
			isdecimal = true
		}
		if isminus {
			l--
		}
	}

	if l <= 0 || !isfloat {
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

func MoneyFormat(n interface{}) (string, error) {
	var buff bytes.Buffer

	e := numberSeparator(n, &buff, ".", 3)
	if e != nil {
		return "", e
	}

	return buff.String(), nil
}

func IDR(n interface{}) (string, error) {
	var buff bytes.Buffer

	buff.Write([]byte("Rp "))

	s, e := MoneyFormat(n)
	if e != nil {
		return "", e
	}

	buff.WriteString(s)

	if strings.Index(s, ",") < 0 {
		buff.Write([]byte(",00"))
	}

	return buff.String(), nil
}
