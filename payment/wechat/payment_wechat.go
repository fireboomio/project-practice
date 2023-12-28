package customize

import (
	"context"
	"crypto/rsa"
	latest "custom-go/generated_latest"
	"custom-go/pkg/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	wxCore "github.com/wechatpay-apiv3/wechatpay-go/core"
	wxVerifiers "github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	wxDecryptors "github.com/wechatpay-apiv3/wechatpay-go/core/cipher/decryptors"
	wxEncryptors "github.com/wechatpay-apiv3/wechatpay-go/core/cipher/encryptors"
	wxDownloader "github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	wxNotify "github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	wxOption "github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	wxUtils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"time"
)

const (
	appName                      string = ""         // 应用名称
	wxAppid                      string = ""         // 应用ID
	wxMiniAppid                  string = ""         // 应用ID
	wxMchID                      string = ""         // 商户号
	wxMchCertificateSerialNumber string = ""         // 商户证书序列号
	wxMchAPIv3Key                string = ""         // 商户APIv3密钥
	wxApiclientKeyFilepath       string = ""         // 商户私钥存放路径
)

var (
	wxMchPrivateKey *rsa.PrivateKey
	wxNotifyHandler *wxNotify.Handler
)

const (
	WxPay PayType = "wxPay"
	PC    string  = "pc"
	H5    string  = "h5"
	MINI  string  = "mini"
)

func init() {
	initWxConfig()

	PayMap[WxPay] = &payAction{
		buildClient: func(ctx context.Context) (client any, err error) {
			client, err = BuildWxPayClient(ctx)
			return
		},
		prepayForApp: func(client any, grc *base.GraphqlRequestContext, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string, mode string) (resp any, err error) {
			if mode == PC {
				resp, err = nativePay(client, totalFee, notifyURI, outTradeNo, productName, expireAt)
			} else if mode == H5 {
				resp, err = jsapiPay(grc, wxAppid, client, totalFee, notifyURI, outTradeNo, productName, expireAt, mode)
			} else if mode == MINI {
				resp, err = jsapiPay(grc, wxMiniAppid, client, totalFee, notifyURI, outTradeNo, productName, expireAt, mode)
			} else {
				err = fmt.Errorf("【WeChat】未知的支付模式: %s", mode)
			}
			return
		},
		PayNotify: func(request *base.ClientRequest) (paymentUpdateInput PaymentUpdateI, err error) {
			paymentUpdateInput, err = PayNotify(request)
			return
		},
		statusQuery: func(client any, i *base.InternalClient, outTradeNo string) (res any, err error) {
			mode, err := getPaymentByOutTradeNo(i, outTradeNo)
			if mode == PC {
				res, err = nativeQuery(client, outTradeNo)
			} else if mode == H5 || mode == MINI {
				res, err = jsapiQuery(client, outTradeNo)
			} else {
				err = fmt.Errorf("【WeChat】未知的支付模式: %s", mode)
			}
			return
		},
		ModifyRequestForPayNotify: func(body *plugins.HttpTransportBody) (*base.ClientRequest, error) {
			return body.Request, nil
		},
	}
}

func getPaymentByOutTradeNo(i *base.InternalClient, outTradeNo string) (mode string, err error) {
	paymentRes, err := latest.Payment__GetOnePaymentByOutTradeNo.Execute(paymentByOrderNumberGetI{OrderNumber: outTradeNo}, i.Context)
	if err != nil {
		return
	}
	if paymentRes.Data.Id == "" {
		err = fmt.Errorf("【WeChat】未找到订单: %s", outTradeNo)
		return
	}
	mode = paymentRes.Data.Mode
	return
}

func getCurrentUserOpenId(grc *base.GraphqlRequestContext, mode string) (openId string, err error) {
	wxs, ok := grc.User.CustomClaims["wxs"]
	if !ok {
		err = errors.New("未获取到微信登录信息")
		return
	}

	var userWxs []*userWx
	wxsBytes, _ := json.Marshal(wxs)
	_ = json.Unmarshal(wxsBytes, &userWxs)
	if len(userWxs) == 0 {
		err = errors.New("未获取到微信登录信息")
		return
	}
	// 根据mode获取对应的openid
	for _, item := range userWxs {
		if item.Platform == mode {
			openId = item.Openid
			break
		}
	}
	if openId == "" {
		err = errors.New("未获取到对应端的微信openid")
		return
	}
	return
}

// jsapiPay 获取prepay_id
func jsapiPay(grc *base.GraphqlRequestContext, appid string, client any, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string, mode string) (resp any, err error) {
	// client
	wxClient := client.(*wxCore.Client)
	svc := jsapi.JsapiApiService{Client: wxClient}
	ctx := context.Background()

	// 过期时间
	timeExpire, err := time.Parse(utils.ISO8601Layout, expireAt)

	// 获取openid
	openid, err := getCurrentUserOpenId(grc, mode)
	if err != nil {
		return
	}

	// 调用微信jsapi接口
	resp, _, err = svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       wxCore.String(appid),
			Mchid:       wxCore.String(wxMchID),
			OutTradeNo:  wxCore.String(outTradeNo),
			Description: wxCore.String(utils.JoinString("-", appName, productName)),
			Attach:      wxCore.String("自定义数据说明"),
			NotifyUrl:   wxCore.String(notifyURI),
			TimeExpire:  &timeExpire,
			Amount: &jsapi.Amount{
				Total: wxCore.Int64(totalFee),
			},
			Payer: &jsapi.Payer{
				Openid: wxCore.String(openid),
			},
		},
	)
	if err != nil {
		return
	}

	prepayResp, ok := resp.(*jsapi.PrepayWithRequestPaymentResponse)
	if !ok {
		log.Println(resp)
	}
	jsonResp, err := json.Marshal(prepayResp)
	if err != nil {
		log.Println("Error marshalling PrepayWithRequestPaymentResponse to JSON:", err)
		return
	}
	resp = string(jsonResp)
	return
}

// nativePay 获取二维码
func nativePay(client any, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string) (resp any, err error) {
	// client
	wxClient := client.(*wxCore.Client)
	svc := native.NativeApiService{Client: wxClient}
	ctx := context.Background()

	// 过期时间
	timeExpire, err := time.Parse(utils.ISO8601Layout, expireAt)

	// 调用微信native接口
	resp, _, err = svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:       wxCore.String(wxAppid),
			Mchid:       wxCore.String(wxMchID),
			OutTradeNo:  wxCore.String(outTradeNo),
			Description: wxCore.String(utils.JoinString("-", appName, productName)),
			Attach:      wxCore.String("自定义数据说明"),
			NotifyUrl:   wxCore.String(notifyURI),
			TimeExpire:  &timeExpire,
			Amount: &native.Amount{
				Total: wxCore.Int64(totalFee),
			},
		},
	)
	if err != nil {
		return
	}

	prepayResp, ok := resp.(*native.PrepayResponse)
	if !ok {
		log.Println(resp)
	}
	resp = *prepayResp.CodeUrl
	return
}

func jsapiQuery(client any, outTradeNo string) (res any, err error) {
	wxClient := client.(*wxCore.Client)
	svc := jsapi.JsapiApiService{Client: wxClient}

	res, _, err = svc.QueryOrderByOutTradeNo(context.Background(),
		jsapi.QueryOrderByOutTradeNoRequest{
			OutTradeNo: wxCore.String(outTradeNo),
			Mchid:      wxCore.String(wxMchID),
		},
	)
	if err != nil {
		return
	}
	return
}

func nativeQuery(client any, outTradeNo string) (res any, err error) {
	wxClient := client.(*wxCore.Client)
	svc := native.NativeApiService{Client: wxClient}

	res, _, err = svc.QueryOrderByOutTradeNo(context.Background(),
		native.QueryOrderByOutTradeNoRequest{
			OutTradeNo: wxCore.String(outTradeNo),
			Mchid:      wxCore.String(wxMchID),
		},
	)
	if err != nil {
		return
	}
	return
}

func initWxConfig() {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	wxMchPrivateKey, _ = wxUtils.LoadPrivateKeyWithPath(wxApiclientKeyFilepath)

	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	_ = wxDownloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, wxMchPrivateKey, wxMchCertificateSerialNumber, wxMchID, wxMchAPIv3Key)
	// 2. 获取商户号对应的WeChat支付平台证书访问器
	certificateVisitor := wxDownloader.MgrInstance().GetCertificateVisitor(wxMchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	wxNotifyHandler, _ = wxNotify.NewRSANotifyHandler(wxMchAPIv3Key, wxVerifiers.NewSHA256WithRSAVerifier(certificateVisitor))
}

// BuildWxPayClient 使用商户私钥等初始化 client，并使它具有自动定时获取WeChat支付平台证书的能力
func BuildWxPayClient(ctx context.Context) (client any, err error) {
	// 使用商户私钥等初始化 client，并使它具有自动定时获取WeChat支付平台证书的能力
	opts := []wxCore.ClientOption{
		wxOption.WithWechatPayAutoAuthCipher(wxMchID, wxMchCertificateSerialNumber, wxMchPrivateKey, wxMchAPIv3Key),
		wxOption.WithWechatPayCipher(
			wxEncryptors.NewWechatPayEncryptor(wxDownloader.MgrInstance().GetCertificateVisitor(wxMchID)),
			wxDecryptors.NewWechatPayDecryptor(wxMchPrivateKey),
		),
	}
	client, err = wxCore.NewClient(ctx, opts...)
	return
}

func PayNotify(request *base.ClientRequest) (paymentUpdateInput PaymentUpdateI, err error) {
	transaction := new(payments.Transaction)
	_, err = wxNotifyHandler.ParseNotifyRequest(context.Background(), request.NewRequest(), transaction)
	if err != nil {
		return
	}

	mappedStatus, ok := StatusMapping[*transaction.TradeState]
	if !ok {
		err = fmt.Errorf("【WeChat】未知的交易状态: %s", *transaction.TradeState)
		return
	}

	paymentUpdateInput = PaymentUpdateI{
		OrderNumber:   *transaction.OutTradeNo,
		PaymentDate:   *transaction.SuccessTime,
		PaymentStatus: mappedStatus,
	}
	return
}
