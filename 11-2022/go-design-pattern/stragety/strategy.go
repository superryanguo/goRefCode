package main

import "fmt"

//策略模式这个定义乍一看起来，还是挺抽象、挺难懂的，这里说的算法并不是我们想找工作准备面试时每天要刷的那种算法；定义一类算法族中的算法族说的要完成的某项任务的归类，举个例子来说比如用户支付，就是个任务类。

//算法族中的每个算法（即策略）则是说的完成这项任务的具体方式，结合我们的例子来说就是可以用支付宝也可以用微信支付这两种方式 (算法) ，来完成我们定义的用户支付这项任务 (算法族)。

//策略模式主要用于允许我们的程序在运行时动态更改一个任务的处理逻辑，常见的应用场景有针对软件用户群体的不同策略切换（用一个烂大街的词儿表达就是千人千面）和业务流程兜底切换。

//注意：这里是为了大家好理解举了支付这个例子，实际上运行时切换支付方式还是挺复杂的，实践的时候你可以先从运行时切换通知用户的任务练起。

//策略模式要解决的问题是，让使用客户端跟具体执行任务的策略解耦，不管使用哪种策略完成任务，不需要更改客户端使用策略的方式。
type PayBehavior interface {
	OrderPay(px *PayCtx)
}

// 具体支付策略实现
// 微信支付
type WxPay struct{}

func (*WxPay) OrderPay(px *PayCtx) {
	fmt.Printf("Wx支付加工支付请求 %v\n", px.payParams)
	fmt.Println("正在使用Wx支付进行支付")
}

// 三方支付
type ThirdPay struct{}

func (*ThirdPay) OrderPay(px *PayCtx) {
	fmt.Printf("三方支付加工支付请求 %v\n", px.payParams)
	fmt.Println("正在使用三方支付进行支付")
}

//有了策略的实现后，还得有个上下文来协调它们，以及持有完成这个任务所必需的那些入参payParams

type PayCtx struct {
	// 提供支付能力的接口实现
	payBehavior PayBehavior
	// 支付参数
	payParams map[string]interface{}
}

func (px *PayCtx) setPayBehavior(p PayBehavior) {
	px.payBehavior = p
}

func (px *PayCtx) Pay() {
	px.payBehavior.OrderPay(px)
}

func NewPayCtx(p PayBehavior) *PayCtx {
	// 支付参数，Mock数据
	params := map[string]interface{}{
		"appId": "234fdfdngj4",
		"mchId": 123456,
	}
	return &PayCtx{
		payBehavior: p,
		payParams:   params,
	}
}

//所有这些代码都准备好后，我们就可以试着运行程序调用它们了。

func main() {
	wxPay := &WxPay{}
	px := NewPayCtx(wxPay)
	px.Pay()
	// 假设现在发现微信支付没钱，改用三方支付进行支付
	thPay := &ThirdPay{}
	px.setPayBehavior(thPay)
	px.Pay()
}
