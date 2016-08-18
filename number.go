package number

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func NumberSeparator(n interface{}, sep string, t int) (string, error) {
	var now_n string
	var temp_n float64
	var isfloat bool
	var isminus bool
	var isdecimal bool

	temp_n, isfloat, err := NumberToFloat64(n)

	switch val := n.(type) {
	case int:
		now_n = strconv.FormatInt(int64(n.(int)), 10)
		temp_n = float64(n.(int))
		isfloat = false
	case int64:
		now_n = strconv.FormatInt(n.(int64), 10)
		temp_n = float64(n.(int64))
		isfloat = false
	case float32:
		now_n = strconv.FormatFloat(float64(n.(float32)), 'f', -1, 32)
		temp_n = float64(n.(float32))
		isfloat = true
	case float64:
		now_n = strconv.FormatFloat(n.(float64), 'f', -1, 64)
		temp_n = n.(float64)
		isfloat = true
	default:
		return "", errors.New(fmt.Sprintln("[Number Separator] Wrong Data Type", val))
	}

	startOffset := 0
	var buff bytes.Buffer
	if temp_n < 0 {
		isminus = true
		startOffset = 1
		buff.WriteByte('-')
	}
	var l int
	if isfloat {
		l = strings.Index(now_n, ".") + 1
		if l > 0 {
			isdecimal = true
		}
		if isminus {
			l--
		}
	}
	if l <= 0 || !isfloat {
		l = len(now_n)
	}

	commaIndex := t - ((l - startOffset) % t)
	if commaIndex == t {
		commaIndex = 0
	}
	var i int
	for i = startOffset; i < l; i++ {
		if commaIndex == t {
			buff.WriteString(sep)
			commaIndex = 0
		}
		commaIndex++
		buff.WriteByte(now_n[i])
	}
	if isdecimal {
		buff.Write([]byte(now_n[i:]))
	}

	return buff.String(), nil
}

func StringRuneSeparator(s, r string, d int) string {

}

func NumberToFloat64(n interface{}) (float64, bool, error) {
	var f bool
	var t float64

	switch val := n.(type) {
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
	var t float64

	switch val := n.(type) {
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

func ProperFloat64(n float64) float64 {
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.15f", n), 64)
	return t
}
