package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// ContainString check if string slice contain value
func ContainString(slice []string, val string) bool {
	for i := range slice {
		if slice[i] == val {
			return true
		}
	}
	return false
}

// DescribeProcess show list searchable field by table name
func DescribeProcess(schema *Schema, arrCommandStr []string) error {
	if len(arrCommandStr) < 2 {
		return errors.New("Required table name")
	}
	tableName := arrCommandStr[1]
	if isOk := ContainString(globalConfig.SupportedSchema, tableName); !isOk {
		fmt.Println("Unsupported table")
		fmt.Println("List supported table:")
		fmt.Println(strings.Join(globalConfig.SupportedSchema, ", "))
		return nil
	}
	fmt.Println("List searchable field:")
	for _, v := range schema.SearchField[tableName] {
		fmt.Println(v)
	}
	fmt.Println("")
	return nil
}

// SearchProcess search record by filter
func SearchProcess(schema *Schema, arrCommandStr []string) ([]DataMap, error) {
	if len(arrCommandStr) < 3 {
		return nil, errors.New("Required for 2 arguments")
	}
	tableName := arrCommandStr[1]
	if isOk := ContainString(globalConfig.SupportedSchema, tableName); !isOk {
		fmt.Println("Unsupported table")
		fmt.Println("List supported table:")
		fmt.Println(strings.Join(globalConfig.SupportedSchema, ", "))
		return nil, nil
	}
	params := strings.Split(arrCommandStr[2], "=")
	if len(params) != 2 {
		return nil, errors.New("Filter format is invalid: require fieldName=fieldValue")
	}
	result := schema.Search(tableName, params[0], params[1])
	return result, nil
}

//RenderTabularView render record by tabular view
func RenderTabularView(tableName string, result []DataMap) {
	data := [][]string{}

	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{}
	rowKey := []string{}
	for _, v := range globalConfig.DisplayField[tableName] {
		headers = append(headers, v.Label)
		rowKey = append(rowKey, v.Key)
	}
	table.SetHeader(headers)
	for i := range result {
		rowValue := []string{}
		item := result[i]
		for _, key := range rowKey {
			if item[key] != nil && reflect.TypeOf(item[key]).Kind() != reflect.Slice {
				rowValue = append(rowValue, fmt.Sprintf("%v", item[key]))
			} else if item[key] != nil && reflect.TypeOf(item[key]).Kind() == reflect.Slice {
				joinS := []string{}
				if slice, castOk := item[key].([]interface{}); castOk {
					for _, v := range slice {
						joinS = append(joinS, fmt.Sprintf("%v", v))
					}
					rowValue = append(rowValue, strings.Join(joinS, ", "))
					continue
				}
				rowValue = append(rowValue, "")
			} else {
				rowValue = append(rowValue, "")
			}
		}
		data = append(data, rowValue)
	}
	table.AppendBulk(data)
	table.Render()
}
