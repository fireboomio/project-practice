package customize

import (
	"context"
	"custom-go/generated"
	"custom-go/pkg/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/cast"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"os"
	"time"
)

type (
	getPendingPaymentsI  = any
	getPendingPaymentsRD = generated.Payment__GetPendingPaymentsResponseData
	cancelOnePaymentI    = generated.Payment__CancelOnePaymentInternalInput
	cancelOnePaymentRD   = generated.Payment__CancelOnePaymentResponseData
	paymentConfI         = generated.Payment__GetPaymentConfInternalInput
	paymentConfRD        = generated.Payment__GetPaymentConfResponseData
)

var (
	getPendingPaymentsQueryPath  = generated.Payment__GetPendingPayments
	cancelOnePaymentMutationPath = generated.Payment__CancelOnePayment
	paymentConfQueryPath         = generated.Payment__GetPaymentConf
)

var (
	StatusMapping = map[string]generated.Freetalk_PaymentStatus{
		"TRADE_FINISHED": generated.Freetalk_PaymentStatus_PAID,
		"TRADE_SUCCESS":  generated.Freetalk_PaymentStatus_PAID,
		"SUCCESS":        generated.Freetalk_PaymentStatus_PAID,
		"TRADE_CLOSED":   generated.Freetalk_PaymentStatus_CANCELLED,
		"WAIT_BUYER_PAY": generated.Freetalk_PaymentStatus_PENDING,
	}
)

func init() {
	base.AddRegisteredHook(startPaymentCron)
}

const (
	paymentConfCronIntervalSecKey = "paymentConf.cronIntervalSec"
	paymentConfStartQueryMinKey   = "paymentConf.startQueryMin"
)

func startPaymentCron(logger echo.Logger) {
	internalClient := plugins.DefaultInternalClient

	cronIntervalSec := os.Getenv(paymentConfCronIntervalSecKey)
	startQueryMin := os.Getenv(paymentConfStartQueryMinKey)
	if len(cronIntervalSec) == 0 {
		cronIntervalSec = "5"
	}
	if len(startQueryMin) == 0 {
		startQueryMin = "0"
	}

	for range time.Tick(time.Duration(cast.ToInt(cronIntervalSec)) * time.Second) {
		paymentsRD, err := getPendingPayments(internalClient)
		if err != nil {
			logger.Errorf("获取待支付订单失败：%v", err)
			continue
		}

		for _, data := range paymentsRD.Data {
			createAt, err := time.Parse(utils.ISO8601Layout, data.CreatedAt)
			if err != nil {
				logger.Errorf("解析订单'%s'的创建时间'%s'失败：%v", data.OrderNumber, data.CreatedAt, err)
				continue
			}
			expireAt, err := time.Parse(utils.ISO8601Layout, data.ExpireAt)
			if err != nil {
				logger.Errorf("解析订单'%s'的过期时间'%s'失败：%v", data.OrderNumber, data.ExpireAt, err)
				continue
			}

			// 创建时间距离现在不足2分钟，暂不查询
			if time.Since(createAt) < time.Duration(cast.ToInt(startQueryMin))*time.Minute {
				logger.Infof("订单未到查询时间: %s", data.OrderNumber)
				continue
			}

			// 订单过期，取消订单
			if time.Since(expireAt) > 0 {
				_, err := cancelOnePayment(internalClient, data.Id)
				if err != nil {
					logger.Errorf("取消订单'%s'失败：%v", data.OrderNumber, err)
				}
				continue
			}

			logger.Infof("开始查询订单: %s", data.OrderNumber)
			// 查询订单状态
			pay, ok := PayMap[PayType(data.PayType)]
			if !ok {
				logger.Errorf("不支持的支付类型：%s", data.PayType)
				continue
			}

			payClient, err := pay.buildClient(context.Background())
			if err != nil {
				logger.Errorf("创建支付客户端失败，订单'%s'：%v", data.OrderNumber, err)
				continue
			}

			result, err := pay.statusQuery(payClient, internalClient, data.OrderNumber)
			if err != nil {
				logger.Errorf("查询订单'%s'的状态失败：%v", data.OrderNumber, err)
				continue
			}

			// 获取订单状态
			var tradeStatus string
			switch data.PayType {
			case "aliPay":
				if resp, ok := result.(*alipay.TradeQueryRsp); ok {
					tradeStatus = string(resp.TradeStatus)
				} else {
					fmt.Println("类型断言失败：result 不是 *TradeQueryRsp 类型")
				}
			case "wxPay":
				if resp, ok := result.(*payments.Transaction); ok {
					if resp.TradeState != nil {
						tradeStatus = *resp.TradeState
					} else {
						logger.Infof("TradeState 字段为空")
					}
				} else {
					logger.Infof("类型断言失败：result 不是 *Transaction 类型")
				}
			default:
				logger.Infof("未知的支付类型")
			}

			logger.Infof("查询到订单'%s'的状态：%v", data.OrderNumber, tradeStatus)

			// 更新订单状态
			orderStatus, ok := StatusMapping[tradeStatus]
			if !ok {
				continue
			}

			contentBytes, err := json.Marshal(result)
			if err != nil {
				logger.Errorf("JSON编码失败：%s", err)
				continue
			}
			paymentResp := string(contentBytes)

			paymentUpdateInput := PaymentUpdateI{
				OrderNumber:   data.OrderNumber,
				PaymentDate:   utils.CurrentDateTime(),
				PaymentStatus: orderStatus,
				PaymentResp:   paymentResp,
			}
			updateResp, err := UpdateOnePayment(internalClient, paymentUpdateInput)
			if err != nil {
				logger.Errorf("更新订单'%s'的状态失败：%v", data.OrderNumber, err)
				continue
			}
			if updateResp.Data.Id == "" {
				logger.Errorf("订单'%s'尚未支付成功：%v", data.OrderNumber, err)
				continue
			}

			orderCalculate, ok := OrderCalculateMap[UnifiedOrderProduct(data.Usage)]
			if !ok {
				logger.Errorf("不支持的产品：%s", data.Usage)
				continue
			}

			duration := orderCalculate.FetchPaymentDuration(updateResp)

			_, err = DurationCreateByUsage(internalClient, string(data.Usage), data.AccountId, duration, "")
			if err != nil {
				logger.Errorf("记录和更新订单'%s'的时长失败：%v", data.OrderNumber, err)
			}
		}
	}
}

// getPendingPayments 获取状态为PENDING的订单
func getPendingPayments(internalClient *base.InternalClient) (data getPendingPaymentsRD, err error) {
	data, err = plugins.ExecuteInternalRequestQueries[getPendingPaymentsI, getPendingPaymentsRD](internalClient, getPendingPaymentsQueryPath, nil)
	return
}

// UpdateOnePayment 完成订单
func UpdateOnePayment(internalClient *base.InternalClient, paymentUpdateInput PaymentUpdateI) (PaymentUpdateRD, error) {
	if paymentUpdateInput.PaymentStatus == AliPaid {
		return plugins.ExecuteInternalRequestMutations[PaymentUpdateI, PaymentUpdateRD](internalClient, paymentUpdatePath, paymentUpdateInput)
	}
	return PaymentUpdateRD{}, nil
}

// cancelOnePayment 取消订单
func cancelOnePayment(internalClient *base.InternalClient, id string) (cancelOnePaymentRD, error) {
	cancelOnePaymentInput := cancelOnePaymentI{Id: id}
	return plugins.ExecuteInternalRequestMutations[cancelOnePaymentI, cancelOnePaymentRD](internalClient, cancelOnePaymentMutationPath, cancelOnePaymentInput)
}
