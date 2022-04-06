package util

import uuid "github.com/satori/go.uuid"

// GetUuid 生成唯一标识UUID
func GetUuid() string {
	return uuid.NewV4().String()
}
