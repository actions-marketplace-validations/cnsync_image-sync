package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var (
	mirrorCtx []string
	finalTags []string
)

type list struct {
	Repository string   `json:"Repository"`
	Tags       []string `json:"Tags"`
}

func main() {
	body := httpclient("https://raw.githubusercontent.com/cnsync/image-sync/main/mirrors.txt")

	mirrorCtx = strings.Split(body, "\n")

	for _, cmd := range mirrorCtx {
		srcRepo, srcTags := listTags(cmd)
		if srcRepo == "" {
			log.Println("Empty tags for command:", cmd)
			return
		}

		srcRe, destRe := ImageContains(srcRepo, "docker.io/cnxyz")

		_, destTag := listTags(destRe)

		if destTag != nil {
			finalTags = removeDuplicates(srcTags, destTag)

			// 使用 sort.Slice 对数组进行排序
			sort.Slice(finalTags, func(i, j int) bool {
				return finalTags[i] < finalTags[j]
			})
		}

		copyTags(srcRe, destRe)
	}

}

func listTags(image string) (string, []string) {
	var out bytes.Buffer

	if image != "" {
		cmd := exec.Command("skopeo", "list-tags", "docker://"+image)
		log.Println("Cmd", cmd.Args)
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("exec.Command 命令执行错误....")
			return "", nil
		}

		var l list
		err = json.Unmarshal(out.Bytes(), &l)
		if err != nil {
			log.Println("json.Unmarshal 转换错误....")
			log.Fatal(err)
		}

		return l.Repository, l.Tags
	}

	return "", nil
}

// ImageContains 镜像名称处理
func ImageContains(repo string, name string) (src, dest string) {
	beginIndex := strings.Index(repo, "/")
	b1 := repo[beginIndex+1:]
	b2 := strings.Replace(b1, "/", "-", -1)
	return repo, name + "/" + b2
}

// httpclient http客户端
func httpclient(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("error:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error:", err)
		return ""
	}
	return string(body)
}

func removeDuplicates(left []string, right []string) []string {
	set := make(map[string]bool)
	for _, r := range right {
		set[r] = true
	}

	var result []string
	for _, l := range left {
		if !set[l] {
			result = append(result, l)
		}
	}
	return result
}

func copyTags(srcRe, destRe string) {

	for _, tag := range finalTags {

		cmd := exec.Command("skopeo", "copy", "--insecure-policy", "--src-tls-verify=false", "--dest-tls-verify=false", "-q", "docker://"+srcRe+":"+tag, "docker://"+destRe+":"+tag)
		log.Printf("CMD:[%s]\n", cmd.Args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Error executing command %s: %s\n", cmd.Args, err)
		}
	}

}
