package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
)

func encode(buf *bytes.Buffer, v reflect.Value, space int) error {
	var ws string
	for i := 0; i < space; i++ {
		ws += " "
	}
	if v.IsZero() {
		return nil
	}
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Bool:
		if v.Bool() == true {
			fmt.Fprintf(buf, "t")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(buf, "#C(%.1f, %.1f)", real(v.Complex()), imag(v.Complex()))
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem(), 0)
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString("\n " + ws)
			}
			if err := encode(buf, v.Index(i), 0); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i), len(v.Type().Field(i).Name)+3); err != nil {
				return err
			}
			if i != v.NumField()-1 {
				buf.WriteString(")\n")
			} else {
				buf.WriteByte(')')
			}
		}
		buf.WriteString(")\n")
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key, 0); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key), 0); err != nil {
				return err
			}
			if i != len(v.MapKeys())-1 {
				buf.WriteString(")\n" + ws)
			} else {
				buf.WriteByte(')')
			}
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Satire          bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
	Float           float64
	Complex         complex128
}

func main() {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Satire:   true,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
		Float:   1.1,
		Complex: 1 + 2i,
	}
	b, err := Marshal(strangelove)
	if err != nil {
		log.Fatal(fmt.Errorf("Marshaling is failed:%v", err))
	}
	os.Stdout.Write(b)
}
