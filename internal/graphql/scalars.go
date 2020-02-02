package graphql

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalStringMap writes the StringMap to GraphQL response.
func MarshalStringMap(m map[string]string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(m)
	})
}

// UnmarshalStringMap returns the value of map[string]string from interface{}.
func UnmarshalStringMap(v interface{}) (map[string]string, error) {
	mi, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%T is not a map", v)
	}

	m := make(map[string]string, len(mi))
	for k, vl := range mi {
		s, ok := vl.(string)
		if !ok {
			return nil, fmt.Errorf("%T is not a string", vl)
		}
		m[k] = s
	}

	return m, nil
}
