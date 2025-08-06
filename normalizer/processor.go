package normalizer

import (
	"fmt"
)

func Process(
	input any,
	fields []FieldConfig,
) (any, error) {
	for _, field := range fields {
		switch {
		case field.Name != "":
			// Modifier-type: work on an existing field value
			// if err := applyModifier(input, field); err != nil {
			// 	return err
			// }
			fmt.Println("Apply modifier: ", field.Name)

		case field.Format != "":
			// Constructor-type: create new value from other fields
			// if err := applyConstructor(input, field); err != nil {
			// 	return err
			// }
			fmt.Println("Apply constructor: ", field.Format)

		default:
			return nil, fmt.Errorf("unsupported field config: %+v", field)
		}
	}

	return nil, fmt.Errorf("Error")
}

type ExpandedField struct {
	ResolvedPath string
	IndexMap     map[string]int
}
