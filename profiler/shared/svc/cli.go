package svc

import "fmt"

func DoCliProcess(args []string) {

	usage := `
Usage:
	app cli arg1
	`

	if len(args) < 2 {
		fmt.Println(usage)
	} else {
		switch args[1] {
		case "cmd_one":
			fmt.Println("This is command no.1")
		case "cmd_two":
			fmt.Println("This is command no.2")
		default:
			fmt.Println(usage)
		}
	}
}
