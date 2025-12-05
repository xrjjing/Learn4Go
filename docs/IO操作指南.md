# Go IO 与文件操作速览

## 目录
- 示例目录：`examples/io/`（bufio_rw / copy_file / zip_basic / json_xml）

## 读取与写入
- bufio 扫描：`bufio.Scanner` 逐行/逐词；适合文本处理。示例 `bufio_rw`。
- 文件拷贝：`io.Copy` 简洁高效；注意关闭文件与权限。示例 `copy_file`。

## 压缩
- `archive/zip`：创建、写入、读取 zip；示例 `zip_basic` 展示内存写入再读取。

## 序列化
- JSON：`encoding/json`，struct tag 控制字段；示例 `json_xml`。
- XML：`encoding/xml`，同样支持 tag。

## 小贴士
- `defer file.Close()` 保证资源释放；大文件循环读写需检查错误。
- 读取二进制建议使用 `bufio.Reader`/`Writer` 提升效率。
- 压缩/解压注意文件大小与目录遍历安全。
