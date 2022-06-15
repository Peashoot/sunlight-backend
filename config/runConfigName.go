package config

import (
	"strings"

	"github.com/google/uuid"
)

// RunConfigName 运行配置名称
type RunConfigName = string

var (
	RCN_AuthSignKey    RunConfigName = "用户认证签名密钥" // 用户认证签名密钥
	RCN_AuthEncryptKey RunConfigName = "用户认证加密密钥" // 用户认证加密密钥

	RCN_UserTokenCacheExpiration   RunConfigName = "登录Token有效期" // 登录Token有效期（单位：分钟）
	RCN_UserCacheExpiration        RunConfigName = "用户缓存有效期"    // 用户缓存有效期（单位：分钟）
	RCN_CategoryCacheExpiration    RunConfigName = "类别缓存有效期"    // 类别缓存有效期（单位：分钟）
	RCN_GroupCacheExpiration       RunConfigName = "群组缓存有效期"    // 群组缓存有效期（单位：分钟）
	RCN_GroupMemberCacheExpiration RunConfigName = "组员缓存有效期"    // 组员缓存有效期（单位：分钟）
	RCN_UserLabelCacheExpiration   RunConfigName = "标签缓存有效期"    // 标签缓存有效期（单位：分钟）
	RCN_UserLoginLockExpiration    RunConfigName = "用户登陆锁有效期"   // 用户登陆锁有效期（单位：秒）
	RCN_UserRegisterLockExpiration RunConfigName = "用户注册锁有效期"   // 用户注册锁有效期（单位：秒）
	RCN_UserModifyLockExpiration   RunConfigName = "用户修改锁有效期"   // 用户修改锁有效期（单位：秒）
	RCN_UserCancelLockExpiration   RunConfigName = "用户注销锁有效期"   // 用户注销锁有效期（单位：秒）
	RCN_LabelAddLockExpiration     RunConfigName = "标签新增锁有效期"   // 标签新增锁有效期（单位：秒）
	RCN_LableWaitLockTimeOut       RunConfigName = "标签等待锁超时时间"  // 标签等待锁超时时间（单位：毫秒）
	RCN_ValidUserCacheExpiration   RunConfigName = "无效用户缓存有效期"  // 无效用户缓存有效期（单位：小时）

	RCN_AliYunOssAccessKeyId        RunConfigName = "阿里云AK"          // 阿里云AK
	RCN_AliYunOssAccessKeySecret    RunConfigName = "阿里云SK"          // 阿里云SK
	RCN_AliYunOssRAMRoleArn         RunConfigName = "阿里云Arn"         // 阿里云Arn
	RCN_AliYunOssRAMRoleSessionName RunConfigName = "阿里云SessionName" // 阿里云SessionName
	RCN_AliYunOssEndPoint           RunConfigName = "阿里云端点"          // 阿里云端点
	RCN_AliYunOssRegionId           RunConfigName = "阿里云区域ID"        // 阿里云区域ID
	RCN_AliYunOssBucketName         RunConfigName = "阿里云OSS桶名"       // 阿里云OSS桶名
	RCN_AliYunConfigUseCName        RunConfigName = "阿里云使用CName"     // 阿里云使用CName
	RCN_AliYunConfigEnableCRC       RunConfigName = "阿里云启用CRC校验"     // 阿里云启用CRC校验
	RCN_AliYunConfigConnTimeout     RunConfigName = "阿里云连接超时"        // 阿里云连接超时（单位：秒）
	RCN_AliYunConfigRwTimeout       RunConfigName = "阿里云读写超时"        // 阿里云读写超时（单位：秒）
)

func Load() {
	Register(RCN_AuthSignKey, strings.ReplaceAll(uuid.NewString(), "-", ""))
	Register(RCN_AuthEncryptKey, strings.ReplaceAll(uuid.NewString(), "-", ""))
	Register(RCN_UserTokenCacheExpiration, 60)
	Register(RCN_UserCacheExpiration, 60*24*14)
	Register(RCN_CategoryCacheExpiration, 60*24*14)
	Register(RCN_GroupCacheExpiration, 60*24*14)
	Register(RCN_GroupMemberCacheExpiration, 60*24*14)
	Register(RCN_UserLabelCacheExpiration, 60*24*14)
	Register(RCN_UserLoginLockExpiration, 30)
	Register(RCN_UserRegisterLockExpiration, 30)
	Register(RCN_UserModifyLockExpiration, 30)
	Register(RCN_UserCancelLockExpiration, 30)
	Register(RCN_LabelAddLockExpiration, 30)
	Register(RCN_LableWaitLockTimeOut, 1000)
	Register(RCN_ValidUserCacheExpiration, 60*24*14)
	Register(RCN_AliYunOssAccessKeyId, "<your access key id>")
	Register(RCN_AliYunOssAccessKeySecret, "<your access key secret>")
	Register(RCN_AliYunOssRAMRoleArn, "<your ram role arn>")
	Register(RCN_AliYunOssRAMRoleSessionName, "<your ram role arn>")
	Register(RCN_AliYunOssEndPoint, "<your oss endpoint>")
	Register(RCN_AliYunOssRegionId, "<your oss region id>")
	Register(RCN_AliYunOssBucketName, "<your oss bucket name>")
	Register(RCN_AliYunConfigUseCName, false)
	Register(RCN_AliYunConfigEnableCRC, false)
	Register(RCN_AliYunConfigConnTimeout, int64(5))
	Register(RCN_AliYunConfigRwTimeout, int64(120))
}
