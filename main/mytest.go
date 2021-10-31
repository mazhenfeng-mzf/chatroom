package main

import (
	"fmt"
)



type Run struct {
	name string
}

func (this *Run) run() {
	fmt.Println(this.name + "is running....")
}




type Human struct {
	name string
	age int
	run *Run
}

xiaoming := Human{
	name: "xiaoming",
	age: 15,
	
}

xiaohong := Human{
	name: "xiaohong",
	age: 16,
}


func main() {

	humanList := make([]*Human, 2, 2)

	append(humanList, &xiaoming)
	append(humanList, &xiaohong)

	for _, human := range humanList {
		fmt.Println(human.name)
		fmt.Println(human.age)

	}



}