package factory

import "fmt"

// 工厂模式（Factory Pattern）是面向对象编程中的常用模式。
// 在 Go 项目开发中，你可以通过使用多种不同的工厂模式，来使代码更简洁明了。
// Go 中的结构体，可以理解为面向对象编程中的类，例如 Person 结构体（类）实现了 Greet 方法。

// 简单工厂模式

type Person1 struct {
	Name string
	Age  int
}

func (p Person1) Greet1() {
	fmt.Printf("Hi! My name is %s", p.Name)
}

func NewPerson1(name string, age int) *Person1 {
	return &Person1{
		Name: name,
		Age:  age,
	}
}

// 抽象工厂模式

type Person2 interface {
	Greet2()
}

type person2 struct {
	Name string
	Age  int
}

func (p person2) Greet2() {
	fmt.Printf("Hi! My name is %s", p.Name)
}

func NewPerson2(name string, age int) Person2 {
	return person2{
		Name: name,
		Age:  age,
	}
}

// 工厂方法模式

type Person3 struct {
	Name string
	Age  int
}

func NewPerson3(age int) func(name string) Person3 {
	fn := func(name string) Person3 {
		return Person3{
			Name: name,
			Age:  age,
		}
	}
	return fn
}
