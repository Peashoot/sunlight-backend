package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/do/responses"
	"github.com/peashoot/sunlight/utils"
)

var (
	sigKey        = []byte("signature_hmac_secret_shared_key") // 签名密钥
	encKey        = []byte("GCM_AES_256_secret_shared_key_32") // 加密密钥
	JwtSigner     *jwt.Signer                                  // JWT签名器
	JwtVerifier   *jwt.Verifier                                // JWT校验器
	AuthNeedParty router.Party                                 // 需要鉴权的分组
)

// AuthClaims token中包含的信息
type AuthClaims struct {
	Code           string `json:"code"`           // 用户编号
	Account        string `json:"account"`        // 用户账号
	LoginTimestamp int64  `json:"loginTimestamp"` // 登录时间戳
}

// Validate 验证Token是否有效
func (claims *AuthClaims) Validate(ctx *context.Context) error {
	token := JwtVerifier.RequestToken(ctx)
	// 判断Redis中是否存在，不存在就认为JWT过期，需要重新登录；单点登录限制，一个用户只允许一个地方登录（这里严重依赖Redis）
	if existToken, err := utils.RedisGet[string](utils.JwtRedisCachePrefix + claims.Code); err != nil {
		return err
	} else if existToken != token {
		return errors.New("unvalid token")
	}
	return nil
}

// Init 初始化JWT鉴权
func Init(app *iris.Application) {
	sigKey = []byte(config.GetValue[string](config.RCN_AuthSignKey))
	encKey = []byte(config.GetValue[string](config.RCN_AuthEncryptKey))
	JwtSigner = jwt.NewSigner(jwt.HS256, sigKey,
		time.Duration(config.GetValue[int](config.RCN_UserTokenCacheExpiration))*time.Minute)
	JwtSigner.WithEncryption(encKey, nil)

	JwtVerifier = jwt.NewVerifier(jwt.HS256, sigKey)
	JwtVerifier.WithDefaultBlocklist()
	JwtVerifier.WithDecryption(encKey, nil)
	// 认证失败的特殊返回修改为JSON
	JwtVerifier.ErrorHandler = func(ctx *context.Context, err error) {
		backDo := responses.NewPackagedRespModel()
		backDo.Code = responses.UnauthorizedCode
		backDo.Msg = "请先登录"
		ctx.StopWithJSON(http.StatusUnauthorized, backDo)
	}
	verifyMiddleware := JwtVerifier.Verify(func() interface{} {
		return new(AuthClaims)
	})
	AuthNeedParty = app.Party("/", verifyMiddleware)
	AuthNeedParty.Use(verifyMiddleware)
}
