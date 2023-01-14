package oth

import (
	"fmt"

	"github.com/jinzhu/copier"
)

type Aoo struct {
	UserName string `copier:"UserName"`
	Age      int
}

type Boo struct {
	NickName string `copier:"UserName"`
	MyAge    int    `copier:"Age"`
}

func Copy() {
	aoo := Aoo{UserName: "lulu", Age: 18}
	boo := Boo{}
	if err := copier.Copy(&boo, &aoo); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", boo)
}
