



<p align="center">
  <a href=""><img src="https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG155.jpg" width="320" height="150" alt="one-api logo"></a>
</p>

<div align="center">

# ALI PAY

_✨ App支付 ✨_

</div>



## 整体交互流程图

![Alt text](https://ft-dev.oss-cn-shanghai.aliyuncs.com/123.jpeg)

以下对重点步骤做简要说明：
• 第 1 步用户在商户 App 客户端/小程序中购买商品下单。
• 第 2 步商户订单信息由商户 App 客户端/小程序发送到服务端。
• 第 3 步商家服务端调用 alipay.trade.app.pay（app支付接口2.0）通过支付宝服务端 SDK 获取 orderStr（orderStr 中包含了订单信息和签名）。
• 第 4 步商家将 orderStr 发送给商户 App 客户端/小程序。
• 第 5 步商家在客户端/小程序发起请求，将 orderStr 发送给支付宝。
• 第 6 步进行支付预下单：支付宝客户端将会按照商家客户端提供的请求参数进行支付预下单。正常场景下，会唤起支付宝收银台等待用户核身；异常场景下，会返回异常信息。
• 第 11 步返回商家 App/小程序：用户在支付宝 App 完成支付后，会跳转回商家页面，并返回最终的支付结果（即同步通知），可查看 同步通知参数说明。
• 第 13 步支付结果异步通知，支付宝会根据步骤3 传入的异步通知地址 notify_url，发送异步通知，可查看 异步通知参数说明。
除了正向支付流程外，支付宝也提供交易查询、关闭、退款、退款查询以及对账等配套 API。

## 功能描述
```go
// 1. 系统和环境要求,需要安装的依赖库：
import (
    "github.com/smartwalle/alipay/v3"
)

// 2.1 配置支付宝客户端
const (
	aliAppid                    string = ""
	alipayPrivateKeyFilePath    string = ""
	alipayAppPublicCertFilePath string = ""
	alipayRootCertFilePath      string = ""
	alipayPublicCertFilePath    string = ""
)
// 2.2 使用 initAlipayConfig() 函数从文件中读取支付宝私钥。
func initAlipayConfig() {
	alipayPrivateKey, _ := os.ReadFile(alipayPrivateKeyFilePath)
	aliPrivateKey = string(alipayPrivateKey)
}

// 2.3 使用 loadCerts() 函数加载支付宝应用公钥证书、支付宝根证书和支付宝公钥证书。
func loadCerts(client *alipay.Client) error {
	_ := client.LoadAppPublicCertFromFile(alipayAppPublicCertFilePath)
	_ := client.LoadAliPayRootCertFromFile(alipayRootCertFilePath)
	_ := client.LoadAliPayPublicCertFromFile(alipayPublicCertFilePath)
	return nil
}

// 3. buildClient 初始化支付宝客户端
func buildClient: func(ctx context.Context) (client any, err error) {}

// 4. prepayForApp 用于处理应用内支付（App支付）的预支付请求，返回拉起 支付宝App 的 sn 串
func prepayForApp: func(client any, grc *base.GraphqlRequestContext, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string, mode string) (resp any, err error) {}

// 5. PayNotify 函数来处理支付宝的异步通知，验证通知的签名，解析并处理通知中的订单号、支付日期和交易状态
func PayNotify: func(request *base.ClientRequest) (paymentUpdateInput PaymentUpdateI, err error) {}

// 6. statusQuery 查询支付订单的状态。使用 aliClient.TradeQuery 根据外部订单号查询订单状态。
func statusQuery: func(client any, i *base.InternalClient, outTradeNo string) (res any, err error){} {}

// 7. 修改原始的支付通知请求
func ModifyRequestForPayNotify: func(body *plugins.HttpTransportBody) (*base.ClientRequest, error) {}

```
