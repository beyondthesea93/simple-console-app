package main

import "testing"

var mockSchema = Schema{
	DataHS: map[string]map[interface{}]DataMap{
		"test": {
			"1": {
				"_id":  1,
				"name": "Linh Nguyen",
				"languages": []interface{}{
					"Swift", "Python", "Golang",
				},
				"alias": "linhnguyen",
			},
			"2": {
				"_id":  2,
				"name": "Vu Nguyen",
				"languages": []interface{}{
					"Swift", "Javascript", "React",
				},
				"alias": "vunguyen",
			},
			"3": {
				"_id":  3,
				"name": "Thanh Nguyen",
				"languages": []interface{}{
					"Java", "Javascript", "React",
				},
				"alias": "thanhnguyen",
			},
		},
	},
	SearchField: map[string][]string{
		"test": []string{
			"_id", "name", "languages", "alias",
		},
	},
}

func TestSearchSchemaSingleValue(t *testing.T) {
	result := mockSchema.Search("test", "alias", "linhnguyen")
	if len(result) != 1 {
		t.Errorf("Unexpected search record, expect 1 record but get %d", len(result))
	}
}

func TestSearchSchemaMultipleValue(t *testing.T) {
	result := mockSchema.Search("test", "languages", "React")
	if len(result) != 2 {
		t.Errorf("Unexpected search record, expect 2 record but get %d", len(result))
	}
}
