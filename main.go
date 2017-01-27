package main

import "fmt"

func main() {
	fmt.Println(ParseSql("SELECT user,name,foo,bar FROM google"))
}
