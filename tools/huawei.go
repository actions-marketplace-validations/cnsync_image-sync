// 批量处理华为云镜像 是否为公共仓库 调用官方API
package main

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/region"
	"time"
)

type T struct {
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Size         int64     `json:"size"`
	IsPublic     bool      `json:"is_public"`
	NumImages    int       `json:"num_images"`
	NumDownload  int       `json:"num_download"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Logo         string    `json:"logo"`
	Url          string    `json:"url"`
	Path         string    `json:"path"`
	InternalPath string    `json:"internal_path"`
	DomainName   string    `json:"domain_name"`
	Namespace    string    `json:"namespace"`
	Tags         []string  `json:"tags"`
	Status       bool      `json:"status"`
	TotalRange   int       `json:"total_range"`
}

var (
	ak = ""
	sk = ""
)

func main() {

}

// UpdateRepo 更新镜像仓库的概要信息
func UpdateRepo() {
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())

	request := &model.UpdateRepoRequest{}
	request.Namespace = "cnxyz"
	request.Repository = "distroless-base"
	request.Body = &model.UpdateRepoRequestBody{
		IsPublic: true,
	}
	response, err := client.UpdateRepo(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}

// ListReposDetails 查询镜像仓库列表
func ListReposDetails() {

	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf("cn-east-3")).
			WithCredential(auth).
			Build())

	request := &model.ListReposDetailsRequest{}
	response, err := client.ListReposDetails(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}
