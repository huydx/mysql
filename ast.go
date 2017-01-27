package main

import "fmt"

type SelectStmt struct {
	Fields    []string
	TableName string
}

func (s SelectStmt) String() string {
	return fmt.Sprintf("Fields: %v, TableName: %s", s.Fields, s.TableName)
}
