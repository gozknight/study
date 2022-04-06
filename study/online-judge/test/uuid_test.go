package test

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	id2 := uuid.NewV4().String()
	id3 := uuid.NewV4().String()
	id4 := uuid.NewV4().String()
	println(id2)
	println(id3)
	println(id4, " ", len(id4))
}
