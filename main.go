package main

import (
	"fmt"
	scp_parser "scp-parser/src"
)

func main() {
	// data := []string{"/scp-1122", "/scp-001", "/scp-002"}
	arraySCP := scp_parser.ParseGetListSCP()[0:2]
	for _, item := range arraySCP {
		unit := scp_parser.ParseGetCurrentSCP(item)
		fmt.Println(unit)
	}
	// res := parseGetCurrentSCP(data[0])
	// fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
