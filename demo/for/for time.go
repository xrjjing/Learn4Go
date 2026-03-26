package main

import (
	"fmt"
	"time"
)

func main(){
	fmt.Println(time.Now())
	fmt.Println(time.DateOnly)
	fmt.Println(time.Date(2026,3,24,10,37,15,122,time.Local))

	var x int = 5
	var y int = 5
	fmt.Println(printForAdd(x,y))
	fmt.Println(printForSub(x,y))
	fmt.Println(printForX(x,y))
	fmt.Println(printForDiv(x,y))
	printForOneToNine()
}

func printForAdd(x,y int) int{
	return  x + y
}

func printForSub(x,y int) int{
	return  x - y
}

func printForX(x,y int) int{
	return  x * y
}

func printForDiv(x,y int) int{
	return  x / y
}

func printForOneToNine(){
	for i := range 10 {
		fmt.Print(i)
	}
}
