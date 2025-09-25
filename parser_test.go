package cvimodelgo

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, _ := os.Open("model.cvimodel")
	defer f.Close()

	yolo, err := ParseModelFile(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", yolo)
}
