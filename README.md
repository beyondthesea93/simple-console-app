# simple-console-app
Simple Console Application

# Requirement
* Go version go1.13.5 windows/amd64

# How To Start Program
* Clone project to directory covered by $GO_PATH

* Install 3rd library:

$ go get "github.com/olekukonko/tablewriter"

$ go get "github.com/hokaccha/go-prettyjson"

$ go get "github.com/json-iterator/go"

* Run program

$ go run . (for Window)

$ go run *.go (for Mac)

* Run Unit test

$ go test -timeout 30s simple-console-app

# List command:
  - search [table-name] [filter-field]=[filter-value] : search record from table, response by JSON.
    - Ex: "search organization tags=Cherry"
  - table [table-name] [filter-field]=[filter-value] : search record from table, return tabular view.
    - Ex: "table user name=Cross Barlow"
  - describe [table-name] : return list searchable field by table
    - Ex: "describe ticket"
