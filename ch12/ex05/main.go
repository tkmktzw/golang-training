package main

import (
	"bytes"
	"encoding/json"
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
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")
	case reflect.Bool:
		if v.Bool() == true {
			fmt.Fprintf(buf, "true")
		} else {
			fmt.Fprintf(buf, "false")
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
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString(",\n " + ws)
			}
			if err := encode(buf, v.Index(i), 0); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case reflect.Struct:
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "\"%s\":", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i), len(v.Type().Field(i).Name)+3); err != nil {
				return err
			}
			if i != v.NumField()-1 {
				buf.WriteString(",\n")
			}
		}
		buf.WriteString("}\n")
	case reflect.Map:
		buf.WriteByte('{')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			//buf.WriteByte('[')
			if err := encode(buf, key, 0); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encode(buf, v.MapIndex(key), 0); err != nil {
				return err
			}
			if i != len(v.MapKeys())-1 {
				buf.WriteString(",\n" + ws)
			}
		}
		buf.WriteByte('}')
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
	Title    string            `json:"Title"`
	Subtitle string            `json:"Subtitle"`
	Year     int               `json:"Year"`
	Color    bool              `json:"Color"`
	Satire   bool              `json:"Satire"`
	Actor    map[string]string `json:"Actor"`
	Oscars   []string          `json:"Oscars"`
	Sequel   *string           `json:"Sequal"`
	Float    float64           `json:"Float"`
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
		Float: 1.1,
	}
	b, err := Marshal(strangelove)
	if err != nil {
		log.Fatal(fmt.Errorf("Marshaling failed:%v", err))
	}
	os.Stdout.Write(b)
	var movie Movie
	if err := json.Unmarshal(b, &movie); err != nil {
		log.Fatal(fmt.Errorf("Unmarshaling failed:%v", err))
	}
	b, err = Marshal(movie)
	if err != nil {
		log.Fatal(fmt.Errorf("Marshaling failed:%v", err))
	}
}
