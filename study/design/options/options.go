package options

import (
	"time"
)

// 选项模式（Options Pattern）也是 Go 项目开发中经常使用到的模式，
// 例如，grpc/grpc-go 的NewServer函数，uber-go/zap 包的New函数都用到了选项模式。
// 使用选项模式，我们可以创建一个带有默认值的 struct 变量，并选择性地修改其中一些参数的值。

const (
	defaultTimeout = 10
	defaultCaching = false
)

type Connection struct {
	addr    string
	cache   bool
	timeout time.Duration
}

type ConnectionOptions struct {
	Caching bool
	Timeout time.Duration
}

func NewDefaultOptions() *ConnectionOptions {

	return &ConnectionOptions{
		Caching: defaultCaching,
		Timeout: defaultTimeout,
	}
}

// NewConnect creates a connection with options.
func NewConnect(addr string, opts *ConnectionOptions) (*Connection, error) {
	return &Connection{
		addr:    addr,
		cache:   opts.Caching,
		timeout: opts.Timeout,
	}, nil
}
