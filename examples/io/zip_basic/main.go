package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// 将内存数据写入 zip 并再读出
func main() {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	f, _ := zw.Create("hello.txt")
	_, _ = f.Write([]byte("hello zip"))
	zw.Close()

	// 保存到文件
	tmp := "zip_demo.zip"
	os.WriteFile(tmp, buf.Bytes(), 0644)
	defer os.Remove(tmp)

	// 读取 zip
	r, err := zip.OpenReader(tmp)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, _ := f.Open()
		data, _ := io.ReadAll(rc)
		rc.Close()
		fmt.Printf("%s => %s\n", f.Name, string(data))
	}
}
