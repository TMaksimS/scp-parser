package main

import (
	"fmt"
	scp_parser "scp-parser/src"
)

func main() {
	scp_parser.CreateTable()
	// data := []string{"/scp-1122", "/scp-001", "/scp-002"}
	arraySCP := scp_parser.ParseGetListSCP()[0:2]
	for _, item := range arraySCP {
		unit := scp_parser.ParseGetCurrentSCP(item)
		fmt.Println(unit)
	}
	// res := scp_parser.ParseGetCurrentSCP(data[0])
	// fmt.Println(res)
	// fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
