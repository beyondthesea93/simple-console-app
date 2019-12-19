package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
)

// DataMap interface map for dynamic object
type DataMap map[string]interface{}

// SliceDataMap slice for data map model
type SliceDataMap []DataMap

// SliceToHashMap convert slice data map to hashmap by key
func (src SliceDataMap) SliceToHashMap(key string) map[interface{}]DataMap {
	resultHashMap := make(map[interface{}]DataMap)
	for i := range src {
		data := src[i]
		resultHashMap[data[key]] = DataMap(data)
	}
	return resultHashMap
}

// Schema schema model
type Schema struct {
	DataHS      map[string]map[interface{}]DataMap
	SearchField map[string][]string
}

// NewSchema custom constructor for schema struct
func NewSchema() *Schema {
	instance := new(Schema)
	instance.DataHS = make(map[string]map[interface{}]DataMap)
	return instance
}

// StreamData stream data from json file
func (src *Schema) StreamData(config *Config) error {
	if config == nil {
		return fmt.Errorf("Missing global config")
	}
	streamChan := make(chan error)
	var mutex sync.Mutex
	for t, p := range globalConfig.DataPath {
		go func(tableNme string, dataPath string) {
			if jsonFile, err := os.Open(dataPath); err != nil {
				streamChan <- err
			} else {
				defer jsonFile.Close()
				byteValue, _ := ioutil.ReadAll(jsonFile)
				data := SliceDataMap([]DataMap{})
				err = JSON.Unmarshal(byteValue, &data)
				if err != nil {
					streamChan <- err
				}

				formattedData := data.SliceToHashMap("_id")
				mutex.Lock()
				defer mutex.Unlock()
				src.DataHS[tableNme] = formattedData
			}
			streamChan <- nil
		}(t, p)
	}
	for range globalConfig.DataPath {
		err := <-streamChan
		if err != nil {
			return err
		}
	}
	src.MergeRelationship()
	src.SearchField = globalConfig.SupportedSearchField
	return nil
}

// MergeRelationship merge relationship by requirement
func (src *Schema) MergeRelationship() {
	for _, v := range src.DataHS["ticket"] {
		if v["organization_id"] != nil {
			if orgID := v["organization_id"]; src.DataHS["organization"][orgID] != nil {
				v["organization_name"] = src.DataHS["organization"][orgID]["name"]
				if src.DataHS["organization"][orgID]["tickets"] == nil {
					src.DataHS["organization"][orgID]["tickets"] = []interface{}{v["subject"]}
					continue
				}
				src.DataHS["organization"][orgID]["tickets"] = append(src.DataHS["organization"][orgID]["tickets"].([]interface{}), v["subject"])
			}
		}
		if v["assignee_id"] != nil {
			if asgID := fmt.Sprintf("%v", v["assignee_id"]); src.DataHS["user"][asgID] != nil {
				if src.DataHS["user"][asgID]["assigned_tickets"] == nil {
					src.DataHS["user"][asgID]["assigned_tickets"] = []interface{}{v["subject"]}
					continue
				}
				src.DataHS["user"][asgID]["assigned_tickets"] = append(src.DataHS["user"][asgID]["assigned_tickets"].([]interface{}), v["subject"])
				v["assignee_name"] = src.DataHS["user"][asgID]["name"]
			}
		}

		if v["submitter_id"] != nil {
			if subID := fmt.Sprintf("%v", v["submitter_id"]); src.DataHS["user"][subID] != nil {
				if src.DataHS["user"][subID]["submitted_tickets"] == nil {
					src.DataHS["user"][subID]["submitted_tickets"] = []interface{}{v["subject"]}
					continue
				}
				src.DataHS["user"][subID]["submitted_tickets"] = append(src.DataHS["user"][subID]["submitted_tickets"].([]interface{}), v["subject"])
				v["submitter_name"] = src.DataHS["user"][subID]["name"]
			}
		}
	}

	for _, v := range src.DataHS["user"] {
		if v["organization_id"] != nil {
			if orgID := v["organization_id"]; src.DataHS["organization"][orgID] != nil {
				v["organization_name"] = src.DataHS["organization"][orgID]["name"]
				if src.DataHS["organization"][orgID]["users"] == nil {
					src.DataHS["organization"][orgID]["users"] = []interface{}{v["name"]}
					continue
				}
				src.DataHS["organization"][orgID]["users"] = append(src.DataHS["organization"][orgID]["users"].([]interface{}), v["name"])
			}
		}
	}
}

// Search query table by fieldName=fieldValue
func (src *Schema) Search(table string, fieldName string, fieldValue interface{}) []DataMap {
	result := []DataMap{}
mainLoop:
	for k, v := range src.DataHS[table] {
		if v[fieldName] != nil {
			if reflect.TypeOf(v[fieldName]).Kind() == reflect.Slice {
				if slice, castOk := v[fieldName].([]interface{}); castOk {
					for _, v := range slice {
						if fmt.Sprintf("%v", v) == fieldValue {
							result = append(result, src.DataHS[table][k])
							continue mainLoop
						}
					}
				}
			} else if fmt.Sprintf("%v", v[fieldName]) == fieldValue {
				result = append(result, src.DataHS[table][k])
			}
		}
	}
	return result
}
