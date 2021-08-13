package database

import (
	"fmt"
	"testing"
)

func TestLink(t *testing.T) {
	db:= Link()
	if db == nil{
		fmt.Println("Link to database failed!")
	}else{
		defer db.Close()
		fmt.Println("Link to database success!")
	}


}