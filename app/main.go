package main

import (
	"fmt"
	"scp-parser/internal/repository/postgresql"
	scpClient "scp-parser/internal/repository/scp-parser"
)

func main() {
	postgresql.CreateTable()
	data := []string{"/scp-1122", "/scp-001", "/scp-002"}
	// arraySCP := scp_parser.ParseGetListSCP()[0:2]
	for _, item := range data {

		unit := scpClient.ParseGetCurrentSCP(item)
		fmt.Println(unit)
	}
	// res := scp_parser.ParseGetCurrentSCP(data[0])
	// fmt.Println(res)
	// fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
