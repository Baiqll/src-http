package httpx

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


func FileServer(w http.ResponseWriter, r *http.Request) {
    // 清理请求路径以防止路径穿越攻击
    file_path := filepath.Clean(r.URL.Path)

    if strings.Contains(file_path, "../") || strings.Contains(file_path, "..\\") {
        http.Error(w, "404 Not Found", http.StatusNotFound)
        return
    }
	
    // // 检查文件是否在允许列表中
    // if _, ok := allowedFiles[fileName]; !ok {
    //     http.Error(w, "403 Forbidden", http.StatusForbidden)
    //     return
    // }

    // 构建文件的完整路径
    file_path = filepath.Join("./", filepath.FromSlash(file_path))

    // 检查文件是否存在
    if file_info, err := os.Stat(file_path); os.IsNotExist(err) || file_info.IsDir() {
        http.Error(w, "404 Not Found", http.StatusNotFound)
        return
    }


    // 读取并提供文件内容
    http.ServeFile(w, r, file_path)
}