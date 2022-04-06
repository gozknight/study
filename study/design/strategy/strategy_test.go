package strategy

import (
	"fmt"
	"testing"
)

func TestStrategy(t *testing.T) {
	operator := Operator{}
	strategy := &add{}
	operator.setStrategy(strategy)
	result := operator.calculate(1, 2)
	fmt.Println("add:", result)

	strategy1 := &reduce{}
	operator.setStrategy(strategy1)
	result = operator.calculate(2, 1)
	fmt.Println("reduce:", result)
}
