package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"

	"github.com/Rxg1898/ProductsStorage/store"
)

var validata = validator.New()

func NewUploader(endpoint, ak, sk string) (store.Uploader, error) {
	uploader := &aliyun{
		Endpoint: endpoint,
		Ak:       ak,
		Sk:       sk,
	}
	if err := uploader.validata(); err != nil {
		return nil, err
	}

	return uploader, nil
}

type aliyun struct {
	Endpoint string `validate:"required,url"`
	Ak       string `validate:"required"`
	Sk       string `validate:"required"`
}

func (a *aliyun) validata() error {
	return validata.Struct(a)
}

func (a *aliyun) UploadFile(bucketName, objectKey, filePath string) error {
	client, err := oss.New(a.Endpoint, a.Ak, a.Sk)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 第一个参数: 上传到oss里面的文件的key(路径), go.sum --->  2021/7/21/go.sum
	// 第二个参数: 需要上传的文件的路径
	err = bucket.PutObjectFromFile(filePath, filePath)
	if err != nil {
		return err
	}

	// 打印下载URL
	// sts, 临时授权token(有效期1天)
	signedURL, err := bucket.SignURL(filePath, oss.HTTPGet, 60*60*24)
	if err != nil {
		return fmt.Errorf("sign file download url error, %s", err)
	}
	fmt.Printf("下载链接: %s\n", signedURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")
	return nil
}
