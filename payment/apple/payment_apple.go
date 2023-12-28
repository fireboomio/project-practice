package customize

import (
	"bytes"
	"context"
	"custom-go/generated"
	"custom-go/pkg/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VerifyRequest struct {
	// Receipt app解析出的票据信息
	Receipt string `json:"receipt-data"`
	// Password App的秘钥
	Password string `json:"password"`
	// ExcludeOldTransactions Set this value to true for the response to include only the latest renewal transaction for any subscriptions. Use this field only for app receipts that contain auto-renewable subscriptions.
	ExcludeOldTransactions bool `json:"exclude-old-transactions"`
}
type VerifyResponse struct {
	Environment        string                `json:"environment"`
	IsRetryable        bool                  `json:"is-retryable"`
	LatestReceipt      string                `json:"latest_receipt,omitempty"`
	LatestReceiptInfo  []*LatestReceiptInfo  `json:"latest_receipt_info,omitempty"`
	PendingRenewalInfo []*PendingRenewalInfo `json:"pending_renewal_info,omitempty"`
	Receipt            *Receipt              `json:"receipt,omitempty"`
	Status             int                   `json:"status"`
}

type LatestReceiptInfo struct {
	CancellationDate              string `json:"cancellation_date"`
	CancellationDateTimestamp     string `json:"cancellation_date_ms"`
	CancellationDatePST           string `json:"cancellation_date_pst"`
	CancellationReason            string `json:"cancellation_reason"`
	ExpiresDate                   string `json:"expires_date"`
	ExpiresDateTimestamp          string `json:"expires_date_ms"`
	ExpiresDatePST                string `json:"expires_date_pst"`
	InAppOwnershipType            string `json:"in_app_ownership_type"`
	IsInIntroOfferPeriod          string `json:"is_in_intro_offer_period"`
	IsTrialPeriod                 string `json:"is_trial_period"`
	IsUpgraded                    string `json:"is_upgraded"`
	OfferCodeRefName              string `json:"offer_code_ref_name"`
	OriginalPurchaseDate          string `json:"original_purchase_date"`
	OriginalPurchaseDateTimestamp string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST       string `json:"original_purchase_date_pst"`
	OriginalTransactionId         string `json:"original_transaction_id"`
	ProductId                     string `json:"product_id"`
	PromotionalOfferId            string `json:"promotional_offer_id"`
	PurchaseDate                  string `json:"purchase_date"`
	PurchaseDateTimestamp         string `json:"purchase_date_ms"`
	PurchaseDatePST               string `json:"purchase_date_pst"`
	Quantity                      string `json:"quantity"`
	SubscriptionGroupIdentifier   string `json:"subscription_group_identifier"`
	WebOrderLineItemId            string `json:"web_order_line_item_id"`
	TransactionId                 string `json:"transaction_id"`
	AppAccountToken               string `json:"app_account_token"`
}

type Receipt struct {
	AdamId                        int64    `json:"adam_id"`
	AppItemId                     int64    `json:"app_item_id"`
	ApplicationVersion            string   `json:"application_version"`
	BundleId                      string   `json:"bundle_id"`
	DownloadId                    int64    `json:"download_id"`
	ExpirationDate                string   `json:"expiration_date"`
	ExpirationDateTimestamp       string   `json:"expiration_date_ms"`
	ExpirationDatePST             string   `json:"expiration_date_pst"`
	InApp                         []*InApp `json:"in_app,omitempty"`
	OriginalApplicationVersion    string   `json:"original_application_version"`
	OriginalPurchaseDate          string   `json:"original_purchase_date"`
	OriginalPurchaseDateTimestamp string   `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST       string   `json:"original_purchase_date_pst"`
	PreorderDate                  string   `json:"preorder_date"`
	PreorderDateTimestamp         string   `json:"preorder_date_ms"`
	PreorderDatePST               string   `json:"preorder_date_pst"`
	ReceiptCreationDate           string   `json:"receipt_creation_date"`
	ReceiptCreationDateTimestamp  string   `json:"receipt_creation_date_ms"`
	ReceiptCreationDatePST        string   `json:"receipt_creation_date_pst"`
	ReceiptType                   string   `json:"receipt_type"`
	RequestDate                   string   `json:"request_date"`
	RequestDateTimestamp          string   `json:"request_date_ms"`
	RequestDatePST                string   `json:"request_date_pst"`
	VersionExternalIdentifier     int64    `json:"version_external_identifier"`
}

type PendingRenewalInfo struct {
	AutoRenewProductId              string `json:"auto_renew_product_id"`
	AutoRenewStatus                 string `json:"auto_renew_status"`
	ExpirationIntent                string `json:"expiration_intent"`
	GracePeriodExpiresDate          string `json:"grace_period_expires_date"`
	GracePeriodExpiresDateTimestamp string `json:"grace_period_expires_date_ms"`
	GracePeriodExpiresDatePST       string `json:"grace_period_expires_date_pst"`
	IsInBillingRetryPeriod          string `json:"is_in_billing_retry_period"`
	OfferCodeRefName                string `json:"offer_code_ref_name"`
	OriginalTransactionId           string `json:"original_transaction_id"`
	PriceConsentStatus              string `json:"price_consent_status"`
	ProductId                       string `json:"product_id"`
	PromotionalOfferId              string `json:"promotional_offer_id"`
}

const (
	Apple             PayType = "applePay"
	AppleSandbox      string  = "https://sandbox.itunes.apple.com/verifyReceipt"
	AppleProd         string  = "https://buy.itunes.apple.com/verifyReceipt"
	ApplePassword     string  = ""
	StatusCodeSandBox         = 21007
)

type InApp struct {
	CancellationDate              string `json:"cancellation_date"`
	CancellationDateTimestamp     string `json:"cancellation_date_ms"`
	CancellationDatePST           string `json:"cancellation_date_pst"`
	CancellationReason            string `json:"cancellation_reason"`
	ExpiresDate                   string `json:"expires_date"`
	ExpiresDateTimestamp          string `json:"expires_date_ms"`
	ExpiresDatePST                string `json:"expires_date_pst"`
	IsInIntroOfferPeriod          string `json:"is_in_intro_offer_period"`
	IsTrialPeriod                 string `json:"is_trial_period"`
	OriginalPurchaseDate          string `json:"original_purchase_date"`
	OriginalPurchaseDateTimestamp string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST       string `json:"original_purchase_date_pst"`
	OriginalTransactionId         string `json:"original_transaction_id"`
	ProductId                     string `json:"product_id"`
	PromotionalOfferId            string `json:"promotional_offer_id"`
	PurchaseDate                  string `json:"purchase_date"`
	PurchaseDateTimestamp         string `json:"purchase_date_ms"`
	PurchaseDatePST               string `json:"purchase_date_pst"`
	Quantity                      string `json:"quantity"`
	TransactionId                 string `json:"transaction_id"`
	WebOrderLineItemId            string `json:"web_order_line_item_id"`
}

func init() {
	initAlipayConfig()

	PayMap[Apple] = &payAction{
		buildClient: func(ctx context.Context) (client any, err error) {
			return "", nil
		},
		prepayForApp: func(client any, grc *base.GraphqlRequestContext, totalFee int64, notifyURI, outTradeNo, productName string, expireAt string, mode string) (resp any, err error) {
			return outTradeNo, nil
		},
		PayNotify: func(request *base.ClientRequest) (paymentUpdateInput PaymentUpdateI, err error) {
			return PaymentUpdateI{}, nil
		},
		statusQuery: func(client any, i *base.InternalClient, outTradeNo string) (res any, err error) {
			return "", nil
		},
		ModifyRequestForPayNotify: func(body *plugins.HttpTransportBody) (*base.ClientRequest, error) {
			return body.Request, nil
		},
	}
}

// finishPayment 完成支付
func finishPayment(grc *base.GraphqlRequestContext, args *AppleVerifyInput, inApps []*InApp) (handledCount int, err error) {
	now := time.Now()
	timeNowFormat, timeNowTimestamp := now.Format(utils.ISO8601Layout), fmt.Sprintf("%v", now.Unix()*1000)
	for _, app := range inApps {
		grc.Logger.Infof("app.TransactionId: %s, app.ExpiresDateTimestamp: %s, timeNowTimestamp: %s", app.TransactionId, app.ExpiresDateTimestamp, timeNowTimestamp)
		if len(app.ExpiresDateTimestamp) > 0 && app.ExpiresDateTimestamp < timeNowTimestamp {
			continue
		}

		paymentByTransactionIdInput := onePaymentByTransactionIdGetI{TransactionId: app.TransactionId}
		paymentByTransactionIdResp, _ := plugins.ExecuteInternalRequestQueries[onePaymentByTransactionIdGetI, onePaymentByTransactionIdGetRD](grc.InternalClient, onePaymentByTransactionIdGetPATH, paymentByTransactionIdInput)
		paymentData := paymentByTransactionIdResp.Data
		var paymentId, productUsage, orderNumber string
		var value int64
		if len(paymentData.Id) > 0 {
			if paymentData.PaymentStatus != generated.Freetalk_PaymentStatus_PENDING {
				grc.Logger.Infof("payment [%s] had %s", paymentData.OrderNumber, paymentData.PaymentStatus)
				continue
			}

			// 获取商品信息
			if _, _, value, err = getProductByProductId(grc, args.AppleProductId); err != nil {
				return
			}

			paymentId, productUsage, orderNumber = paymentData.Id, string(paymentData.Usage), paymentData.OrderNumber
		} else if paymentId, productUsage, value, orderNumber, err = genPayment(grc, args, app.TransactionId); err != nil {
			return
		}

		// 完成订单
		paymentUpdateInput := PaymentUpdateI{
			OrderNumber:   orderNumber,
			PaymentDate:   timeNowFormat,
			PaymentStatus: generated.Freetalk_PaymentStatus_PAID,
			Sn:            args.Receipt,
		}
		if _, err = plugins.ExecuteInternalRequestMutations[PaymentUpdateI, PaymentUpdateRD](grc.InternalClient, paymentUpdatePath, paymentUpdateInput); err != nil {
			return
		}

		// 创建 durationHistory
		if _, err = DurationCreateByUsage(grc.InternalClient, productUsage, args.AccountId, value, paymentId); err != nil {
			return
		}

		handledCount++
	}
	return
}

func genPayment(grc *base.GraphqlRequestContext, args *AppleVerifyInput, transactionId string) (paymentId string, productUsage string, value int64, orderNumber string, err error) {
	// 获取商品信息
	productId, productUsage, value, err := getProductByProductId(grc, args.AppleProductId)
	if err != nil {
		return
	}

	// 下单
	unifiedOrderI := unifiedOrderInput{
		AccountId:     args.AccountId,
		Product:       productUsage,
		ProductId:     productId,
		PayType:       "applePay",
		TransactionId: transactionId,
	}
	_, paymentId, orderNumber, err = UnifiedOrder(grc, &unifiedOrderI)
	if err != nil {
		return
	}
	return
}

// getProductByProductId 根据 appleProductId 获取商品信息
func getProductByProductId(grc *base.GraphqlRequestContext, appleProductId string) (productId string, product string, value int64, err error) {
	membershipByProductIdGetInput := membershipByProductIdGetI{
		ProductId: appleProductId,
	}
	membershipRes, err := plugins.ExecuteInternalRequestQueries[membershipByProductIdGetI, membershipByProductIdGetRD](grc.InternalClient, membershipByProductIdGetPATH, membershipByProductIdGetInput)
	if err != nil {
		return
	}
	if membershipRes.Data.Id == "" {
		durationPackageByProductIdGetInput := durationPackageByProductIdGetI{
			ProductId: appleProductId,
		}
		durationPackageRes, _ := plugins.ExecuteInternalRequestQueries[durationPackageByProductIdGetI, durationPackageByProductIdGetRD](grc.InternalClient, durationPackageByProductIdGetPATH, durationPackageByProductIdGetInput)
		return durationPackageRes.Data.Id, string(DurationPackage), durationPackageRes.Data.Value, nil
	}
	return membershipRes.Data.Id, string(Membership), 0, nil
}

// IosVerify 订单 Receipt 正确性验证
func IosVerify(grc *base.GraphqlRequestContext, args *AppleVerifyInput) (handledCount int, err error) {
	reqData := VerifyRequest{
		Receipt:  args.Receipt,
		Password: ApplePassword,
	}
	request := func(url string) (jsonResp *VerifyResponse, isTest bool, err error) {
		jsonResp, err = requestVerify(url, reqData)
		if err != nil {
			return
		}
		if jsonResp.Status == StatusCodeSandBox {
			isTest = true
			return
		}
		return
	}
	// 生产验证
	JsonResp, isTest, err := request(AppleProd)
	if err != nil {
		return
	}
	if isTest {
		// 沙盒验证
		JsonResp, isTest, err = request(AppleSandbox)
		if err != nil {
			return
		}
		if isTest {
			err = fmt.Errorf("[appStore PayVerify] testurl return should in testmod")
			return
		}
	}
	if JsonResp.Status != 0 {
		err = fmt.Errorf("[appStore PayVerify] JsonResp.Status[%d]!=0", JsonResp.Status)
		return
	}
	handledCount, err = finishPayment(grc, args, JsonResp.Receipt.InApp)
	return
}

func requestVerify(url string, reqData VerifyRequest) (jsonResp *VerifyResponse, err error) {
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(reqData)
	if err != nil {
		return
	}
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Post(url, "application/json", buf)
	// 报错或者没有返回信息
	if err != nil || resp == nil {
		return
	}
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			fmt.Println(e.Error())
			return
		}
	}(resp.Body)
	jsonResp = &VerifyResponse{}
	if resp.StatusCode != http.StatusOK {
		//if resp.StatusCode == http.StatusRequestTimeout { // 比如超时啊什么的 需要让其通过校验
		//	jsonResp.Status = 0
		//}
		return
	}
	jsonRespByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = json.Unmarshal(jsonRespByte, jsonResp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
