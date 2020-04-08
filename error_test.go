package main

import (
	"fmt"
	"testing"
)

func Test_AppError_Error(t *testing.T) {
	err := AppError{"error here"}
	if err.Error() != "error here" {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "err.Error()", "error here", err.Error()))
	}
}
