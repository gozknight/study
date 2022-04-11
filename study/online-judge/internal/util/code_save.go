package util

import "os"

// CodeSave 保存代码
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUuid()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Write(code)
	return path, nil
}
