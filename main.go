package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	mirrorCtx []string
)

const MaxConcurrent = 5 // 最大并发数

func main() {
	body := httpclient("https://raw.githubusercontent.com/DaoCloud/public-image-mirror/main/mirror.txt")

	SplitMirrors(body)

	pool, _ := ants.NewPool(MaxConcurrent)

	defer pool.Release()

	for _, ctx := range mirrorCtx {

		dest, src, _ := ImageContains(listTags(ctx), "docker.io/cnxyz")
		// 组装名称

		log.Println(dest, src)

	}
}

func executeCommand(command string) error {
	// 执行命令的逻辑
	fmt.Println("Executing command:", command)
	return nil
}

// SplitMirrors 切割字符串
func SplitMirrors(body string) {
	temp := strings.Split(body, "\n")

	// 去重
	mirrorCtx = RemoveRepeatedElement(temp)
}

type list struct {
	Repository string   `json:"Repository"`
	Tags       []string `json:"Tags"`
}

func listTags(image string) *list {
	cmd := exec.Command("skopeo", "list-tags", "docker://"+image)
	cmd.Env = os.Environ()
	res, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	var l list

	err = json.Unmarshal(res, &l)
	if err != nil {
		log.Fatal(err)
	}

	return &l
}

// ImageContains 镜像名称处理
func ImageContains(str *list, name string) (src, dest string, tags []string) {

	beginIndex := strings.Index(str.Repository, "/")

	b1 := str.Repository[beginIndex+1:]

	b2 := strings.Replace(b1, "/", "-", -1)

	return str.Repository, name + "/" + b2, str.Tags
}

// httpclient http客户端
func httpclient(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	return string(body)
}

// RemoveRepeatedElement 通过map键的唯一性去重
func RemoveRepeatedElement(s []string) []string {
	result := make([]string, 0)
	m := make(map[string]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
