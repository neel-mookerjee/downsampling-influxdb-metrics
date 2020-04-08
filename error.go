package main

import (
	"fmt"
)

type AppError struct {
	What string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%v", e.What)
}
