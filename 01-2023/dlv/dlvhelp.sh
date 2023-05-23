go install github.com/go-delve/delve/cmd/dlv@latest  

go build -gcflags '-N -l' test.go

dlv debug//命令行进入包所在目录，然后输入 dlv debug 命令进入调试
break main.main
breakpoints
restart重新执行程序

dlv trace
dlv exec

通过vars命令可以查看全部包级的变量
(dlv) vars main

next
continue
dlv trace -ebpf foo

通过stack查看当前执行函数的栈帧信息
进入函数之后可以通过args和locals命令查看函数的参数和局部变量
在执行到main函数断点时，可以disassemble反汇编命令查看main函数对应的汇编



dlv exec test
Type 'help' for list of commands.
(dlv)

首先，使用 ps -ef | grep xxxx 查看卡死的进程的pid，假设卡死的进程pid=1234，然后执行 ./dlv attach 1234 就 可以附着到卡死的进程上了。如果进程是使用非当前用户启动的话，则要加上 sudo 才行

在dlv里，执行 dump ~/sample.core , 就会把当前整个进程的状态都记录下来，那样子就可以下载coredump文件

假设二进制文件名问：sample_bin，则在开发机上执行 dlv core sample_bin sample.core 就能加载到coredump文件进行排查问题了。


    bt，查看栈帧list，查看当前栈帧运行到的代码grs，查看当前所有的goroutine列表gr 4，切换到第四个groutineprint xxx，打印当前栈帧中的变量值args，打印函数参数值locals，打印当前栈帧的local 变量


    优先切换到 goroutine 1 看看是否就是卡死的代码所在
    bt then
    到栈帧5，就是属于我们的业务逻辑代码了。所以此时我执行 frame 5 指令
    list 查看当前栈帧运行到的代码print xxx 查看变量值
    线上调试的时候，想执行 list 查看当前栈帧的代码，出现了类似 Command failed: open /path/to/the/mainfile.go: no such file or directory 的错误。即使下载了coredump文件和二进制到自己的开发机上进行调试，也会出现这个问题，因为二进制是另一个同事编译的，Go的编译默认会带上代码的完整路径，而同事的GOPATH和我是不一样的，这个要怎么解决呢？

终极解决办法，编译参数加上 trimpath：go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH 这样子编译的代码路径，就会去掉你本人的文件夹路径，使用相对于GOPATH的路径了，别人在你编译的二进制上获取的coredump，也可以正常调试了。

如果编译前没有加上这个指令 ，则可以在 dlv 里执行 ： config substitute-path otherPath mypath。

如果这个path弄错了，怎么清除弄错的配置呢？直接执行 config substitute-path otherPath 即可

注意 mypath一定要填绝大路径，使用 ~/ 这样子的路径会有问题。
print变量的时候，如果是字符串变量，好像默认只打印前64个字符，这明显不够，可以通过设置：config max-string-len 99999 让打印的字符串长度加长。此时可以执行 config -save 把这个配置永久保存下来。
