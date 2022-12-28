//建造者模式，也有翻译成生成器模式的，大家看到后知道他们是一个东西，都是Builer Pattern翻译过来的就行。它是一种对象构建模式，是将一个复杂对象的构建与它的表示分离，使得同样的构建过程可以创建不同的表示。 那么什么情况下适合使用建造模式呢？

//当要构建的对象很大并且需要多个步骤时，使用构建器模式，有助于减小构造函数的大小。
//如果你是写过Java程序一定对下面这类代码很熟悉。

//Coffee.builder().name("Latti").price("30").build()
//当然，自己给Coffee类加上构建模式，还是需要写不少额外的代码，得给 Coffee 类加一个静态内部类 CoffeeBuilder，用CoffeeBuilder，去建造Coffee类的对象。
package main

import (
	"fmt"
	"time"
)

type DBPool struct {
	dsn             string
	maxOpenConn     int
	maxIdleConn     int
	maxConnLifeTime time.Duration
}

//need the error, so make it a new type builder?!
type DBPoolBuilder struct {
	DBPool
	err error
}

func Builder() *DBPoolBuilder {
	b := new(DBPoolBuilder)
	// 设置 DBPool 属性的默认值
	b.DBPool.dsn = "127.0.0.1:3306"
	b.DBPool.maxConnLifeTime = 1 * time.Second
	b.DBPool.maxOpenConn = 30
	return b
}

func (b *DBPoolBuilder) DSN(dsn string) *DBPoolBuilder {
	if b.err != nil {
		return b
	}
	if dsn == "" {
		b.err = fmt.Errorf("invalid dsn, current is %s", dsn)
	}

	b.DBPool.dsn = dsn
	return b
}

func (b *DBPoolBuilder) MaxOpenConn(connNum int) *DBPoolBuilder {
	if b.err != nil {
		return b
	}
	if connNum < 1 {
		b.err = fmt.Errorf("invalid MaxOpenConn, current is %d", connNum)
	}

	b.DBPool.maxOpenConn = connNum
	return b
}

func (b *DBPoolBuilder) MaxConnLifeTime(lifeTime time.Duration) *DBPoolBuilder {
	if b.err != nil {
		return b
	}
	if lifeTime < 1*time.Second {
		b.err = fmt.Errorf("connection max life time can not litte than 1 second, current is %v", lifeTime)
	}

	b.DBPool.maxConnLifeTime = lifeTime
	return b
}

func (b *DBPoolBuilder) Build() (*DBPool, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.DBPool.maxOpenConn < b.DBPool.maxIdleConn {
		return nil, fmt.Errorf("max total(%d) cannot < max idle(%d)", b.DBPool.maxOpenConn, b.DBPool.maxIdleConn)
	}
	return &b.DBPool, nil
}
func main() {
	dbPool, err := dbpool.Builder().DSN("localhost:3306").MaxOpenConn(50).MaxConnLifeTime(0 * time.Second).Build()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dbPool)
}

//另外在建造者过程的每个参数步骤里，我们都借用了之前提到的处理 Go Error 的方式，把在外部调用时的错误判断，分散到了每个步骤里。
