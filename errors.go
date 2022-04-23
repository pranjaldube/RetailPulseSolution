package main

import "fmt"

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
