<p align="center">
  <a href=""><img src="https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG156.jpg" width="320" height="140" alt="one-api logo"></a>
</p>

<div align="center">

# WeChat PAY

_✨ PC、小程序、H5 支付 ✨_

</div>

## 整体交互流程图
![Alt text](https://ft-dev.oss-cn-shanghai.aliyuncs.com/126.bmp)


## 功能描述

### 发送请求

先初始化一个 `core.Client` 实例，再向微信支付发送请求。

```go
package main

import (
	"context"
	"log"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/certificates"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func main() {
	var (
		mchID                      string = "190000****"                                // 商户号
		mchCertificateSerialNumber string = "3775B6A45ACD588826D15E583A95F5DD********"  // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"          // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Fatalf("new wechat pay client err:%s", err)
	}
	
	// 发送请求，以下载微信支付平台证书为例
	// https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
	svc := certificates.CertificatesApiService{Client: client}
	resp, result, err := svc.DownloadCertificates(ctx)
	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
}
```

`resp` 是反序列化（UnmarshalJSON）后的应答。上例中是 `services/certificates` 包中的 `*certificates.Certificate`。

`result` 是 `*core.APIResult` 实例，包含了完整的请求报文 `*http.Request` 和应答报文 `*http.Response`。

#### 名词解释

+ **商户 API 证书**，是用来证实商户身份的。证书中包含商户号、证书序列号、证书有效期等信息，由证书授权机构（Certificate Authority ，简称 CA）签发，以防证书被伪造或篡改。如何获取请见 [商户 API 证书](https://wechatpay-api.gitbook.io/wechatpay-api-v3/ren-zheng/zheng-shu#shang-hu-api-zheng-shu) 。
+ **商户 API 私钥**。商户申请商户 API 证书时，会生成商户私钥，并保存在本地证书文件夹的文件 apiclient_key.pem 中。

+ **微信支付平台证书**。微信支付平台证书是指由微信支付负责申请的，包含微信支付平台标识、公钥信息的证书。商户使用微信支付平台证书中的公钥验证应答签名。获取微信支付平台证书需通过 [获取平台证书列表](https://wechatpay-api.gitbook.io/wechatpay-api-v3/ren-zheng/zheng-shu#ping-tai-zheng-shu) 接口下载。
+ **证书序列号**。每个证书都有一个由 CA 颁发的唯一编号，即证书序列号。扩展阅读 [如何查看证书序列号](https://wechatpay-api.gitbook.io/wechatpay-api-v3/chang-jian-wen-ti/zheng-shu-xiang-guan#ru-he-cha-kan-zheng-shu-xu-lie-hao) 。
+ **微信支付 APIv3 密钥**，是在回调通知和微信支付平台证书下载接口中，为加强数据安全，对关键信息 `AES-256-GCM` 加密时使用的对称加密密钥。

## 更多示例

### 以 [JSAPI下单](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_1.shtml) 为例

```go
import (
	"log"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

svc := jsapi.JsapiApiService{Client: client}
// 得到prepay_id，以及调起支付所需的参数和签名
resp, result, err := svc.PrepayWithRequestPayment(ctx,
	jsapi.PrepayRequest{
		Appid:       core.String("wxd678efh567hg6787"),
		Mchid:       core.String("1900009191"),
		Description: core.String("Image形象店-深圳腾大-QQ公仔"),
		OutTradeNo:  core.String("1217752501201407033233368018"),
		Attach:      core.String("自定义数据说明"),
		NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
		Amount: &jsapi.Amount{
			Total: core.Int64(100),
		},
		Payer: &jsapi.Payer{
			Openid: core.String("oUpF8uMuAJO_M2pxb1Q9zNjWeS6o"),
		},
	},
)

if err == nil {
	log.Println(resp)
} else {
	log.Println(err)
}
```

### 以 [查询订单](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml) 为例

```go
import (
	"log"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

svc := jsapi.JsapiApiService{Client: client}

resp, result, err := svc.QueryOrderById(ctx,
	jsapi.QueryOrderByIdRequest{
		TransactionId: core.String("4200000985202103031441826014"),
		Mchid:         core.String("1900009191"),
	},
)

if err == nil {
	log.Println(resp)
} else {
	log.Println(err)
}

```





## 回调通知的验签与解密

1. 使用微信支付平台证书（验签）和商户 APIv3 密钥（解密）初始化 `notify.Handler`
2. 调用 `handler.ParseNotifyRequest` 验签，并解密报文。

### 初始化
+ 方法一（大多数场景）：先手动注册下载器，再获取微信平台证书访问器。

适用场景： 仅需要对回调通知验证签名并解密的场景。例如，基础支付的回调通知。

```go
ctx := context.Background()
// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
err := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIV3Key)
// 2. 获取商户号对应的微信支付平台证书访问器
certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
// 3. 使用证书访问器初始化 `notify.Handler`
handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
```

+ 方法二：像 [发送请求](#发送请求) 那样使用 `WithWechatPayAutoAuthCipher` 初始化 `core.Client`，然后再用client进行接口调用。

适用场景：需要对回调通知验证签名并解密，并且后续需要使用 Client 的场景。例如，电子发票的回调通知，验签与解密后还需要通过 Client 调用用户填写抬头接口。

```go
ctx := context.Background()
// 1. 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
opts := []core.ClientOption{
	option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
}
client, err := core.NewClient(ctx, opts...)	
// 2. 获取商户号对应的微信支付平台证书访问器
certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
// 3. 使用证书访问器初始化 `notify.Handler`
handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
// 4. 使用client进行接口调用
// ...
```

+ 方法三：使用本地的微信支付平台证书和商户 APIv3 密钥初始化 `Handler`。

适用场景：首次通过工具下载平台证书到本地，后续使用本地管理的平台证书进行验签与解密。

```go
// 1. 初始化商户API v3 Key及微信支付平台证书
mchAPIv3Key := "<your apiv3 key>"
wechatPayCert, err := utils.LoadCertificate("<your wechat pay certificate>")
// 2. 使用本地管理的微信支付平台证书获取微信支付平台证书访问器
certificateVisitor := core.NewCertificateMapWithList([]*x509.Certificate{wechatPayCert})
// 3. 使用apiv3 key、证书访问器初始化 `notify.Handler`
handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
```

建议：为了正确使用平台证书下载管理器，你应阅读并理解 [如何使用平台证书下载管理器](FAQ.md#如何使用平台证书下载管理器)。

### 验签与解密

将支付回调通知中的内容，解析为 `payments.Transaction`。

```go
transaction := new(payments.Transaction)
notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, transaction)
// 如果验签未通过，或者解密失败
if err != nil {
	fmt.Println(err)
	return
}
// 处理通知内容
fmt.Println(notifyReq.Summary)
fmt.Println(transaction.TransactionId)
```

将 SDK 未支持的回调消息体，解析至 `map[string]interface{}`。


```go
content := make(map[string]interface{})
notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, &content)
// 如果验签未通过，或者解密失败
if err != nil {
	fmt.Println(err)
	return
}
// 处理通知内容
fmt.Println(notifyReq.Summary)
fmt.Println(content)
```

## 敏感信息加解密

为了保证通信过程中敏感信息字段（如用户的住址、银行卡号、手机号码等）的机密性，

+ 微信支付要求加密上行的敏感信息
+ 微信支付会加密下行的敏感信息


### （推荐）使用敏感信息加解密器

敏感信息加解密器 `cipher.Cipher` 能根据 API 契约自动处理敏感信息：

+ 发起请求时，开发者设置原文，加密器自动加密敏感信息，并设置 `Wechatpay-Serial` 请求头
+ 收到应答时，解密器自动解密敏感信息，开发者得到原文

使用敏感信息加解密器，只需通过 `option.WithWechatPayCipher` 为 `core.Client` 添加加解密器：

```go
client, err := core.NewClient(
    context.Background(),
// 一次性设置 签名/验签/敏感字段加解密，并注册 平台证书下载器，自动定时获取最新的平台证书
    option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
    option.WithWechatPayCipher(
        encryptors.NewWechatPayEncryptor(downloader.MgrInstance().GetCertificateVisitor(mchID)),
        decryptors.NewWechatPayDecryptor(mchPrivateKey),
    ),
)
```
