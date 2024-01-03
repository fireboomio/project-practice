<p align="center">
  <a href=""><img src="https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG157.jpg" width="320" height="140" alt="one-api logo"></a>
</p>

<div align="center">

# Apple PAY

_✨ Apple 支付 ✨_

</div>

## 整体交互流程图
![Alt text](https://ft-dev.oss-cn-shanghai.aliyuncs.com/127.png)


直接购买流程：

1. `客户端` 发起购买，`服务端` 创建订单，返回 `订单编号`

2. `客户端` 拉起苹果支付，获取支付结果收据 `receipt`

3. `客户端` 提交苹果支付收据 `receipt` 到 `服务端` 进行校验



订阅流程：

1. `客户端` 向 `服务端` 发起校验，校验是否可以发起该档位订阅（若已经处于订阅中，则返回不可以购买）

2. 若可以订阅对应档位，`客户端` 发起支付购买, 获取支付收据 `receipt`

3. `客户端` 提交苹果支付收据 `receipt` 到 `服务端` 进行校验，下发，记录订阅数据

4. `服务端` 定时扫描 `周期订阅收据`，产生新的收据做商品的下发

5. `服务端` 收到苹果服务器 `回调` 通知收据，通过原始订单号映射业务账户，`服务端` 进行商品下发和订阅状态更新


## 功能描述

结构体

```go
type VerifyRequest struct {
	Receipt string `json:"receipt-data"`
	Password string `json:"password"`
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
	...
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
	...
}



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
```

配置

```go
const (
	Apple             PayType = "applePay"
	AppleSandbox      string  = "https://sandbox.itunes.apple.com/verifyReceipt"
	AppleProd         string  = "https://buy.itunes.apple.com/verifyReceipt"
	ApplePassword     string  = "02e0e0b64bdc4b17995d096d3b522c19"
	StatusCodeSandBox         = 21007
)1
```

验证逻辑

```go
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

```


服务端定时扫描订阅

```go
func startPaymentAppleSubscription(logger echo.Logger) {
	if !strings.EqualFold(os.Getenv(appleSubscriptionRenewalSwitch), "1") {
		return
	}

	internalClient := plugins.DefaultInternalClient
	renewalFunction := func() {
		now := time.Now()
		manyMembershipResp, _ := plugins.ExecuteInternalRequestQueries[manyMembershipI, manyMembershipRD](internalClient, manyMembershipQueryPath, manyMembershipI{})
		grc := &base.GraphqlRequestContext{Logger: logger, InternalClient: internalClient}
		for _, membershipItem := range manyMembershipResp.Data {
			paymentDateGte := now.Add(timeDay * time.Duration(membershipItem.Lifespan) * -1)
			paymentDateLte := paymentDateGte.Add(timeDay * 1)
			iosRenewalInput := iosRenewalI{PaymentDateGte: paymentDateGte.Format(utils.ISO8601Layout), PaymentDateLte: paymentDateLte.Format(utils.ISO8601Layout)}
			iosRenewalResp, _ := plugins.ExecuteInternalRequestQueries[iosRenewalI, iosRenewalRD](internalClient, iosRenewalQueryPath, iosRenewalInput)
			for _, renewalItem := range iosRenewalResp.Data {
				verifyInput := &AppleVerifyInput{Receipt: renewalItem.Sn, AccountId: renewalItem.AccountId, AppleProductId: membershipItem.ProductId}
				_, _ = IosVerify(grc, verifyInput)
			}
		}
	}
	renewalFunction()
	for range time.Tick(time.Hour) {
		renewalFunction()
	}
}

```