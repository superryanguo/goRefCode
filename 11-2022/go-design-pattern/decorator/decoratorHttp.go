package main

HTTP 相关的一个示例
我们再来看一个处理 HTTP 请求的相关的例子。

先看一个简单的 HTTP Server 的代码。

package main
import (
    "fmt"
    "log"
    "net/http"
    "strings"
)
func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("--->WithServerHeader()")
        w.Header().Set("Server", "HelloServer v0.0.1")
        h(w, r)
    }
}
func hello(w http.ResponseWriter, r *http.Request) {
    log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
    fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}
func main() {
    http.HandleFunc("/v1/hello", WithServerHeader(hello))
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
上面代码中使用到了修饰模式，WithServerHeader() 函数就是一个 Decorator，其传入一个 http.HandlerFunc，然后返回一个改写的版本。上面的例子还是比较简单，用 WithServerHeader() 就可以加入一个 Response 的 Header。

于是，这样的函数我们可以写出好些个。如下所示，有写 HTTP 响应头的，有写认证 Cookie 的，有检查认证Cookie的，有打日志的……

package main
import (
    "fmt"
    "log"
    "net/http"
    "strings"
)
func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("--->WithServerHeader()")
        w.Header().Set("Server", "HelloServer v0.0.1")
        h(w, r)
    }
}
func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("--->WithAuthCookie()")
        cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
        http.SetCookie(w, cookie)
        h(w, r)
    }
}
func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("--->WithBasicAuth()")
        cookie, err := r.Cookie("Auth")
        if err != nil || cookie.Value != "Pass" {
            w.WriteHeader(http.StatusForbidden)
            return
        }
        h(w, r)
    }
}
func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("--->WithDebugLog")
        r.ParseForm()
        log.Println(r.Form)
        log.Println("path", r.URL.Path)
        log.Println("scheme", r.URL.Scheme)
        log.Println(r.Form["url_long"])
        for k, v := range r.Form {
            log.Println("key:", k)
            log.Println("val:", strings.Join(v, ""))
        }
        h(w, r)
    }
}
func hello(w http.ResponseWriter, r *http.Request) {
    log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
    fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}
func main() {
    http.HandleFunc("/v1/hello", WithServerHeader(WithAuthCookie(hello)))
    http.HandleFunc("/v2/hello", WithServerHeader(WithBasicAuth(hello)))
    http.HandleFunc("/v3/hello", WithServerHeader(WithBasicAuth(WithDebugLog(hello))))
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
多个修饰器的 Pipeline
在使用上，需要对函数一层层的套起来，看上去好像不是很好看，如果需要 decorator 比较多的话，代码会比较难看了。嗯，我们可以重构一下。

重构时，我们需要先写一个工具函数——用来遍历并调用各个 decorator：

type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc
func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
    for i := range decors {
        d := decors[len(decors)-1-i] // iterate in reverse
        h = d(h)
    }
    return h
}
然后，我们就可以像下面这样使用了。

http.HandleFunc("/v4/hello", Handler(hello,
                WithServerHeader, WithBasicAuth, WithDebugLog))
这样的代码是不是更易读了一些？pipeline 的功能也就出来了。

泛型的修饰器
不过，对于 Go 的修饰器模式，还有一个小问题 —— 好像无法做到泛型，就像上面那个计算时间的函数一样，其代码耦合了需要被修饰的函数的接口类型，无法做到非常通用，如果这个事解决不了，那么，这个修饰器模式还是有点不好用的。

因为 Go 语言不像 Python 和 Java，Python是动态语言，而 Java 有语言虚拟机，所以他们可以干好些比较变态的事，然而 Go 语言是一个静态的语言，这意味着其类型需要在编译时就要搞定，否则无法编译。不过，Go 语言支持的最大的泛型是 interface{} 还有比较简单的 reflection 机制，在上面做做文章，应该还是可以搞定的。

废话不说，下面是我用 reflection 机制写的一个比较通用的修饰器（为了便于阅读，我删除了出错判断代码）

func Decorator(decoPtr, fn interface{}) (err error) {
    var decoratedFunc, targetFunc reflect.Value
    decoratedFunc = reflect.ValueOf(decoPtr).Elem()
    targetFunc = reflect.ValueOf(fn)
    v := reflect.MakeFunc(targetFunc.Type(),
            func(in []reflect.Value) (out []reflect.Value) {
                fmt.Println("before")
                out = targetFunc.Call(in)
                fmt.Println("after")
                return
            })
    decoratedFunc.Set(v)
    return
}
上面的代码动用了 reflect.MakeFunc() 函数制出了一个新的函数其中的 targetFunc.Call(in) 调用了被修饰的函数。关于 Go 语言的反射机制，推荐官方文章 —— 《The Laws of Reflection》，在这里我不多说了。

上面这个 Decorator() 需要两个参数，

第一个是出参 decoPtr ，就是完成修饰后的函数
第二个是入参 fn ，就是需要修饰的函数
这样写是不是有些二？的确是的。不过，这是我个人在 Go 语言里所能写出来的最好的的代码了。如果你知道更多优雅的，请你一定告诉我！

好的，让我们来看一下使用效果。首先假设我们有两个需要修饰的函数：

func foo(a, b, c int) int {
    fmt.Printf("%d, %d, %d \n", a, b, c)
    return a + b + c
}
func bar(a, b string) string {
    fmt.Printf("%s, %s \n", a, b)
    return a + b
}
然后，我们可以这样做：

type MyFoo func(int, int, int) int
var myfoo MyFoo
Decorator(&myfoo, foo)
myfoo(1, 2, 3)
你会发现，使用 Decorator() 时，还需要先声明一个函数签名，感觉好傻啊。一点都不泛型，不是吗？

嗯。如果你不想声明函数签名，那么你也可以这样

mybar := bar
Decorator(&mybar, bar)
mybar("hello,", "world!")
好吧，看上去不是那么的漂亮，但是 it works。看样子 Go 语言目前本身的特性无法做成像 Java 或 Python 那样，对此，我们只能多求 Go 语言多放糖了！
