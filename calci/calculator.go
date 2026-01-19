package main

import (
	"fmt"
	"errors"
)

func main(){
	var n int
	fmt.Print("Number or operations : ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Error : ", err);
		return
	}

	fmt.Println("Order : num1 operator num2")

	for i := 1; i <= n; i++ {
		var num1, num2, ans float64
		var op string
		_, err = fmt.Scan(&num1, &op, &num2)
		if err != nil {
			fmt.Println("Error : ", err)
			return
		}

		ans, err = calc(num1, num2, op)
		if err != nil {
			fmt.Println("Error : ", err)
			continue
		}
		fmt.Print(ans, "\n\n")
	}
}

func calc(num1, num2 float64, op string) (float64, error) {
	switch op {
		case "+" :
			return num1 + num2, nil
		case "-" : 
			return num1 - num2, nil
		case "*" : 
			return num1 * num2, nil
		case "/" : 
			if num2 == 0{
				return 0, errors.New("Cannot divide with zero\n")
			}
			return num1 / num2, nil
		default : 
			return 0, errors.New("Invalid operator\n")
	}
}