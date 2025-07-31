package cli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

// Parse will fill CLI struct and parse flags. The CLI struct must use `flag` struct tags.
func Parse(target interface{}) {
	val := reflect.ValueOf(target).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		flagTag := field.Tag.Get("flag")
		helpText := field.Tag.Get("help")
		if flagTag == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			ptr := val.Field(i).Addr().Interface().(*string)
			flag.StringVar(ptr, flagTag, *ptr, helpText)
		case reflect.Bool:
			ptr := val.Field(i).Addr().Interface().(*bool)
			flag.BoolVar(ptr, flagTag, *ptr, helpText)
		case reflect.Int:
			ptr := val.Field(i).Addr().Interface().(*int)
			flag.IntVar(ptr, flagTag, *ptr, helpText)
		default:
			fmt.Fprintf(os.Stderr, "⚠️ unsupported CLI flag type: %s\n", field.Type.Name())
		}
	}

	flag.Parse()
}
