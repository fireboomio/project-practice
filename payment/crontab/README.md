<p align="center">
  <a href=""><img src="https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG158.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LTAI5tP8iEBMfrzVNSW7k148%2F20231228%2Foss-cn-shanghai%2Fs3%2Faws4_request&X-Amz-Date=20231228T104253Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3Bfilename%3D%22WechatIMG158.jpg%22&X-Amz-Signature=01e77dfb699d6358635474f2599be4d1d745799a3d26ae8cd65f9e8eae8481ee" width="320" height="140" alt="one-api logo"></a>
</p>

<div align="center">

# Payment Crontab

_✨ 支付定时任务 ✨_

</div>

## 整体交互流程图
![Alt text](https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG39164.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LTAI5tP8iEBMfrzVNSW7k148%2F20231228%2Foss-cn-shanghai%2Fs3%2Faws4_request&X-Amz-Date=20231228T111708Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3Bfilename%3D%22WechatIMG39164.jpg%22&X-Amz-Signature=c53051ef39d67f4aad4ab0e667a5f39270449ee5013640d8765db5c79b7e04b4)


## 功能描述
用于定期检查所有等待支付的订单，并根据支付状态进行相应的处理，包括查询订单状态、取消过期订单以及更新订单状态。

定时查询：定期检查所有等待支付的订单。
订单状态更新：根据支付提供商返回的状态更新订单。
取消过期订单：自动取消已经过期但未支付的订单。
实现步骤

1. 初始化定时任务
在 init 函数中调用 startPaymentCron，启动支付定时任务。

    ```go
    func init() {
        base.AddRegisteredHook(startPaymentCron)
    }
    ```

2. 定时任务执行
    使用 time.Tick 设置定时任务的间隔。
    ```go
    for range time.Tick(time.Duration(cast.ToInt(cronIntervalSec)) * time.Second) {
        // 任务执行逻辑
    }
    ```

3. 处理待支付订单
    对于每个待支付订单，检查订单的创建时间和过期时间。
    如果订单未过期，根据订单的支付类型创建相应的支付客户端。
    使用 statusQuery 方法查询订单的支付状态。
    ```go
    paymentsRD, err := getPendingPayments(internalClient)
    // 处理每个订单
    ```

4. 订单状态更新
    使用 UpdateOnePayment 函数更新订单的支付状态。
    ```go
    paymentUpdateInput := PaymentUpdateI{ /* ... */ }
    updateResp, err := UpdateOnePayment(internalClient, paymentUpdateInput)
    ```

5. 取消过期订单
    如果订单已过期，使用 cancelOnePayment 函数取消订单。
    ```go
    _, err := cancelOnePayment(internalClient, data.Id)
    ```

6. 处理特定的支付类型
    根据订单的支付类型（支付宝或微信支付），处理返回的支付状态。
    ```go
    // 示例：处理支付宝支付类型
    if resp, ok := result.(*alipay.TradeQueryRsp); ok {
        tradeStatus = string(resp.TradeStatus)
    }
    ```
