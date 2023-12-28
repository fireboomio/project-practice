package customize

import (
	"context"
	"custom-go/pkg/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/tidwall/sjson"
	"net/url"
	"os"
	"time"
)

const (
	aliAppid                    string = ""
	alipayPrivateKeyFilePath    string = ""
	alipayAppPublicCertFilePath string = ""
	alipayRootCertFilePath      string = ""
	alipayPublicCertFilePath    string = ""
)

var (
	aliPrivateKey string
	aliClient     *alipay.Client
)

const (
	AliPay           PayType = "aliPay"
	productCode      string  = "QUICK_MSECURITY_PAY"
	outTradeNoField  string  = "out_trade_no"
	paymentDateField string  = "gmt_payment"
	tradeStatusField string  = "trade_status"
	AliPayTimeLayout string  = "2006-01-02 15:04:05"
	AliPaid                  = "PAID"
)

func init() {
	initAlipayConfig()

	PayMap[AliPay] = &payAction{
		buildClient: func(ctx context.Context) (client any, err error) {
			aliClient, err = alipay.New(aliAppid, aliPrivateKey, true)
			if err != nil {
				err = fmt.Errorf("【支付宝】初始化支付宝客户端失败: %s", err)
				return
			}
			err = loadCerts(aliClient)
			if err != nil {
				err = fmt.Errorf("【支付宝】加载证书发生错误: %s", err)
				return
			}
			return aliClient, nil
		},
		prepayForApp: func(client any, grc *base.GraphqlRequestContext, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string, mode string) (resp any, err error) {
			grc.Logger.Infof("【支付宝】notifyUrl: %s, outTradeNo: %s \n", notifyURI, outTradeNo)
			aliClient := client.(*alipay.Client)
			var p = alipay.TradeAppPay{}
			p.NotifyURL = notifyURI
			p.Subject = productName
			p.OutTradeNo = outTradeNo

			// 设置支付宝格式的过期时间
			parse, err := time.Parse(utils.ISO8601Layout, expireAt)
			aliExpireAt := parse.Add(8 * time.Hour).Format(AliPayTimeLayout)
			grc.Logger.Infof("过期时间为：%s", aliExpireAt)
			p.TimeExpire = aliExpireAt

			p.TotalAmount = fmt.Sprintf("%.2f", float64(totalFee)/100)
			p.ProductCode = productCode
			resp, err = aliClient.TradeAppPay(p)
			grc.Logger.Infof("【支付宝】拉起app串：%s\n", resp)
			if err != nil {
				return
			}
			return
		},
		PayNotify: func(request *base.ClientRequest) (paymentUpdateInput PaymentUpdateI, err error) {
			u, err := url.ParseQuery(string(request.OriginBody))
			if err != nil {
				return
			}

			err = aliClient.VerifySign(u)
			if err != nil {
				err = fmt.Errorf("【支付宝】异步通知验证签名发生错误: %s", err)
				return
			}

			outTradeNo := u.Get(outTradeNoField)
			paymentDate := u.Get(paymentDateField)
			tradeStatus := u.Get(tradeStatusField)

			// 处理 tradeStatus
			mappedStatus, ok := StatusMapping[tradeStatus]
			if !ok {
				err = fmt.Errorf("【支付宝】未知的交易状态: %s", tradeStatus)
				return
			}

			// 处理 paymentDate
			parsedPaymentDate, err := time.Parse(AliPayTimeLayout, paymentDate)
			if err != nil {
				err = fmt.Errorf("【支付宝】解析支付日期发生错误: %s", err)
				return
			}
			paymentDateFormat := parsedPaymentDate.Add(-8 * time.Hour).Format(utils.ISO8601Layout)

			paymentUpdateInput = PaymentUpdateI{
				OrderNumber:   outTradeNo,
				PaymentDate:   paymentDateFormat,
				PaymentStatus: mappedStatus,
			}
			return
		},
		statusQuery: func(client any, i *base.InternalClient, outTradeNo string) (res any, err error) {
			var p = alipay.TradeQuery{}
			p.OutTradeNo = outTradeNo
			resp, err := aliClient.TradeQuery(p)
			if err != nil {
				return
			}
			if resp.IsSuccess() == false {
				err = fmt.Errorf(resp.SubMsg)
				return
			}
			res = resp
			return
		},
		ModifyRequestForPayNotify: func(body *plugins.HttpTransportBody) (*base.ClientRequest, error) {
			// 解析原始的URL编码的form数据
			values, err := url.ParseQuery(string(body.Request.OriginBody))
			if err != nil {
				return nil, err
			}

			// 将处理后的数据设置为data字段的值
			modifyBody, err := sjson.Set("{}", "data", values.Encode())
			if err != nil {
				return nil, err
			}

			body.Request.Body = []byte(modifyBody)

			// 修改 "Content-Type" 为 "application/json"
			body.Request.Headers["Content-Type"] = "application/json"

			return body.Request, nil
		},
	}
}

// 读取支付宝私钥
func initAlipayConfig() {
	alipayPrivateKey, err := os.ReadFile(alipayPrivateKeyFilePath)
	if err != nil {
		err = fmt.Errorf("【支付宝】读取支付宝私钥失败: %v", err)
	}
	aliPrivateKey = string(alipayPrivateKey)
}

// 加载证书和内容加密密钥
func loadCerts(client *alipay.Client) error {
	if err := client.LoadAppPublicCertFromFile(alipayAppPublicCertFilePath); err != nil {
		return err
	}
	if err := client.LoadAliPayRootCertFromFile(alipayRootCertFilePath); err != nil {
		return err
	}
	if err := client.LoadAliPayPublicCertFromFile(alipayPublicCertFilePath); err != nil {
		return err
	}
	return nil
}
