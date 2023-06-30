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

type list struct {
	Repository string   `json:"Repository"`
	Tags       []string `json:"Tags"`
}

var FinalTags []string

var HUB string

func init() {
	HUB = os.Getenv("HUB")
}

func main() {
	if HUB == "huawei" {
		body := httpclient("https://raw.githubusercontent.com/cnsync/image-sync/main/mirrors-docker.txt")
		mirrorCtx := strings.Split(body, "\n")
		ExecCommand(mirrorCtx, "swr.cn-east-3.myhuaweicloud.com/cnxyz")
	} else if HUB == "aliyun" {
		body := httpclient("https://raw.githubusercontent.com/cnsync/image-sync/main/mirrors-aliyun.txt")
		mirrorCtx := strings.Split(body, "\n")
		ExecCommand(mirrorCtx, "registry.cn-hangzhou.aliyuncs.com/cnxyz")
	} else {
		body := httpclient("https://raw.githubusercontent.com/cnsync/image-sync/main/mirrors-huawei.txt")
		mirrorCtx := strings.Split(body, "\n")
		ExecCommand(mirrorCtx, "docker.io/cnxyz")
	}

}

func ExecCommand(mirrorCtx []string, hub string) {

	for _, cmd := range mirrorCtx {

		srcRepo, srcTags := listTags(cmd)
		if srcRepo == "" {
			log.Println("Empty tags for command:", cmd)
		}

		srcRe, destRe := ImageContains(srcRepo, hub)

		_, destTag := listTags(destRe)
		//log.Println("-------destTag-----", destTag)

		if destTag != nil {
			FinalTags = removeDuplicates(srcTags, destTag)

			//log.Println("-------FinalTags-----", FinalTags)

			for _, tag := range TagsContains(FinalTags) {
				copyTags(srcRe, destRe, tag)
			}
		} else {

			//log.Println("-----srcTags-------", srcTags)
			for _, tag := range TagsContains(srcTags) {
				copyTags(srcRe, destRe, tag)
			}
		}
	}
}

func TagsContains(strs []string) []string {
	var a []string
	for _, str := range strs {
		if !strings.Contains(str, ".sig") {
			a = append(a, str)
		}
	}

	// 排序
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	return a
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
			log.Println("exec.Command 命令执行错误: ", err)
			return "", nil
		}

		var l list
		err = json.Unmarshal(out.Bytes(), &l)
		if err != nil {
			log.Println("json.Unmarshal 转换错误: ", err)
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
		log.Println("error142:", err)
		return ""
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error148:", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error153:", err)
		return ""
	}
	return string(body)
}

func removeDuplicates(left, right []string) []string {
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

func copyTags(srcRe, destRe, tag string) {
	cmd := exec.Command(
		"skopeo",
		"copy",
		"--insecure-policy",
		"--src-tls-verify=false",
		"--dest-tls-verify=false",
		"-q",
		"docker://"+srcRe+":"+tag,
		"docker://"+destRe+":"+tag)
	log.Println("Cmd", cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error executing command %s: %s\n", cmd.Args, err)
	}
}
