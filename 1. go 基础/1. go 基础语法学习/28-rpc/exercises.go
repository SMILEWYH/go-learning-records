//go:build ignore

package main

// 练习 1：实现远程运费计算。
// 要求：定义 ShippingService.Calculate，参数含 SubtotalCents int64、Express bool，结果含 FeeCents int64。
// 普通运费 600 分，满 10000 分免普通运费；加急额外 1000 分；负金额返回 error。
// 常规：(9999,false)->600；(10000,false)->0；(10000,true)->1000。
// 用独立 rpc.Server 注册，通过真实本机客户端调用；第二个参数必须是结果指针。

// 练习 2：区分本地连接失败和远程业务失败。
// 自行定义 QuoteRequest（Quantity int、UnitPriceCents int64）、QuoteResponse（TotalCents int64）和 PriceService.Quote。
// Quote 拒绝 Quantity<=0，合法时写入数量乘单价的结果；注册自己的本机 RPC 服务。
// 要求：调用 Quote(Quantity=0)，用 errors.As 检查 rpc.ServerError；调用不存在的方法也应报错。
// 再关闭 client，尝试一次新调用，确认失败且不会把之前成功的 reply 当成新结果。
// 每次调用使用独立结果变量；写注释说明服务端的 error 类型为什么不能完整跨网络保留。

// 练习 3：批量异步报价。
// 要求：同一个 client 用 Go 发起三个请求，每个请求有自己的 reply 和容量至少为 1 的 Done。
// 单价 500，数量分别 1、2、3，完成后按输入顺序输出 [500,1000,1500]。
// 边界：将其中一个数量改成 0，收集该项 error，其余成功结果仍保留。
// 不假定网络完成顺序；不在 Done 前读取 reply；为底层连接设置 deadline 并最终 Close。

func main() {
}
