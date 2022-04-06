package single

import "sync"

// 单例模式
// 单例模式（Singleton Pattern），是最简单的一个模式。
// 在 Go 中，单例模式指的是全局只有一个实例，并且它负责创建自己的对象。
// 单例模式不仅有利于减少内存开支，还有减少系统性能开销、防止多个实例产生冲突等优点。

// 饿汉式
type singleton1 struct{}

var ins1 *singleton1 = &singleton1{}

func NewSingleton1() *singleton1 {
	return ins1
}

// 懒汉式-线程不安全
type singleton2 struct{}

var ins2 *singleton2

func NewSingleton2() *singleton2 {
	if ins2 == nil {
		ins2 = &singleton2{}
	}
	return ins2
}

// 懒汉式-线程安全
type singleton3 struct{}

var ins3 *singleton3
var mutex *sync.Mutex

func NewSingleton3() *singleton3 {
	if ins3 == nil {
		mutex.Lock()
		if ins3 == nil {
			ins3 = &singleton3{}
		}
		mutex.Unlock()
	}
	return ins3
}

// Golang自带的更优雅的单例模式实现
type singleton4 struct{}

var ins4 *singleton4
var once *sync.Once

func NewSingleton4() *singleton4 {
	once.Do(func() {
		ins4 = &singleton4{}
	})
	return ins4
}
