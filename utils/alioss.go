package utils

import (
	"io"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/log"
)

// AliOssHelper 阿里OSS帮助类
type AliOssHelper struct {
	Client *oss.Client // 连接客户端
	Bucket *oss.Bucket // OSS存储桶
}

var aliOss *AliOssHelper

func GetAliOss() *AliOssHelper {
	if aliOss == nil {
		aliOss := &AliOssHelper{}
		aliOss.Init()
		return aliOss
	}
	return aliOss
}

// Init 初始化
func (aliOss *AliOssHelper) Init() (err error) {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := config.GetValue[string](config.RCN_AliYunOssEndPoint)
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := config.GetValue[string](config.RCN_AliYunOssAccessKeyId)
	accessKeySecret := config.GetValue[string](config.RCN_AliYunOssAccessKeySecret)
	bucketName := config.GetValue[string](config.RCN_AliYunOssBucketName)
	aliOss.Client, err = oss.New(endpoint, accessKeyId, accessKeySecret,
		oss.UseCname(config.GetValue[bool](config.RCN_AliYunConfigUseCName)),
		oss.SecurityToken("<yourSecurityToken>"),
		oss.Timeout(config.GetValue[int64](config.RCN_AliYunConfigConnTimeout),
			config.GetValue[int64](config.RCN_AliYunConfigRwTimeout)),
		oss.EnableCRC(config.GetValue[bool](config.RCN_AliYunConfigEnableCRC)))
	if err != nil {
		log.Error("[AliOssHelper.Init]", "ali oss init failure, please check params, err:", err)
		return
	}
	var isExists bool
	isExists, err = aliOss.Client.IsBucketExist(bucketName)
	if err != nil {
		log.Error("[AliOssHelper.Init]", "check exist of bucket appear an exception, err:", err)
		return
	}
	if !isExists {
		err = aliOss.Client.CreateBucket(bucketName)
		log.Error("[AliOssHelper.Init]", "bucket doesn't exist, appear an exception while creating bucket, err:", err)
		return
	}
	aliOss.Bucket, err = aliOss.Client.Bucket(bucketName)
	return
}

// UploadFileReader 上传文件流
func (aliOss *AliOssHelper) UploadFileReader(objectName string, fileReader io.Reader) (err error) {
	// 将文件流上传至exampledir目录下的exampleobject.txt文件。
	err = aliOss.Bucket.PutObject(objectName, fileReader)
	if err != nil {
		log.Error("[AliOssHelper.UploadFileReader]", "fail to upload file stream to bucket, the exception is:", err)
	}
	return
}

// UploadFile 上传文件url
func (aliOss *AliOssHelper) UploadFile(objectName, fileURL string) (err error) {
	var res *http.Response
	// 指定待上传的网络流。
	res, err = http.Get(fileURL)
	if err != nil {
		log.Error()
		return
	}
	err = aliOss.Bucket.PutObject(objectName, io.Reader(res.Body))
	if err != nil {
		log.Error("[AliOssHelper.UploadFile]", "fail to upload file url to bucket, the exception is:", err)
	}
	return
}

// GetUploadCredentials 获取上传凭证
func (aliOss *AliOssHelper) GetUploadCredentials() (*sts.Credentials, error) {
	var credentials sts.Credentials
	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	client, err := sts.NewClientWithAccessKey(
		config.GetValue[string](config.RCN_AliYunOssRegionId),
		config.GetValue[string](config.RCN_AliYunOssAccessKeyId),
		config.GetValue[string](config.RCN_AliYunOssAccessKeySecret))
	if err != nil {
		return &credentials, err
	}
	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = config.GetValue[string](config.RCN_AliYunOssRAMRoleArn)
	request.RoleSessionName = config.GetValue[string](config.RCN_AliYunOssRAMRoleSessionName)

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		return &credentials, err
	}
	credentials = response.Credentials
	return &credentials, nil
}
