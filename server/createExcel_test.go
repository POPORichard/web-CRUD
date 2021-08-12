package server

import (
	"sync"
	"testing"
)

func TestWriteToExcel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	WriteToExcel(nil,&wg)
}