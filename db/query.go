package db

import (
	"reflect"
	"strconv"
	"strings"
)

type query struct {
	orderByField string
	orderDir     string
}

func parse(q string) *query {
	res := &query{}
	normalizedStr := strings.ToUpper(q)
	idx := strings.LastIndex(normalizedStr, "ORDER BY")
	if idx >= 0 {
		off := strings.TrimSpace(q[idx+len("ORDER BY"):])
		field_dir := strings.Split(off, " ")
		if len(field_dir) == 1 {
			res.orderByField = field_dir[0]
		} else {
			res.orderByField = field_dir[0]
			res.orderDir = strings.ToUpper(field_dir[1])
		}

	}
	return res
}

type genericJson struct {
	fname string
	data  map[string]interface{}
}

type byCustomField struct {
	array []genericJson
	field string
	asc   bool
}

func (s *byCustomField) Len() int {
	return len(s.array)
}
func (s *byCustomField) Swap(i, j int) {
	s.array[i], s.array[j] = s.array[j], s.array[i]
}
func (s *byCustomField) Less(i, j int) bool {
	a := s.array[i].data[s.field]
	b := s.array[j].data[s.field]

	if a == nil && b == nil {
		return false
	}

	if a == nil && b != nil {
		return false
	}

	if a != nil && b == nil {
		return true
	}

	switch ta := a.(type) {
	case string:
		strAsIntA, strAsIntA_e := strconv.ParseInt(ta, 10, 64)

		switch tb := b.(type) {
		case string:
			strAsIntB, strAsIntB_e := strconv.ParseInt(tb, 10, 64)

			if strAsIntA_e == nil && strAsIntB_e == nil {
				if s.asc {
					return strAsIntA < strAsIntB
				} else {
					return strAsIntA > strAsIntB
				}

			} else {
				if s.asc {
					return ta < tb
				} else {
					return ta > tb
				}

			}

		default:
			panic("not yet implemented: " + reflect.TypeOf(b).String())
		}

	default:
		panic("not yet implemented: " + reflect.TypeOf(a).String())

	}

}
