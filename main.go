package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	mirrorCtx []string
)

const MaxConcurrent = 5 // 最大并发数

func main() {
	body := httpclient("https://raw.githubusercontent.com/DaoCloud/public-image-mirror/main/mirror.txt")

	SplitMirrors(body)

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, MaxConcurrent)

	for _, cmd := range mirrorCtx {
		wg.Add(1)
		semaphore <- struct{}{} // 申请一个信号量，控制并发数量

		go func(command string) {
			defer func() {
				<-semaphore // 释放信号量
				wg.Done()
			}()

			err := executeCommand(command)
			if err != nil {
				fmt.Printf("Error executing command %s: %s\n", command, err)
			}
		}(cmd)
	}

	wg.Wait()
}

func executeCommand(command string) error {

	tags := listTags(command)

	if tags == nil {
		fmt.Println("Empty tags for command:", command)
		return nil
	}

	dest, src, _ := ImageContains(tags, "docker.io/cnxyz")

	// 组装名称
	// log.Println(dest, src)

	// 执行命令的逻辑
	fmt.Println("Executing command:", command)
	cmd := exec.Command("skopeo", "copy", "--insecure-policy", "--src-tls-verify=false", "--dest-tls-verify=false", "-q", "docker://"+src, "docker://"+dest)

	log.Printf("CMD:[%s]\n", cmd.Args)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	// 执行命令的逻辑
	//fmt.Println("Executing command:", command)
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
	var out bytes.Buffer

	cmd := exec.Command("skopeo", "list-tags", "docker://"+image)
	//log.Println("Cmd", cmd.Args)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil
	}

	var l list
	err = json.Unmarshal(out.Bytes(), &l)
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return string(body)
}

// RemoveRepeatedElement 通过map键的唯一性去重
func RemoveRepeatedElement(s []string) []string {
	result := make([]string, 0)
	m := make(map[string]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
