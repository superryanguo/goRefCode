package db

// db.go
type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}

//go get -u github.com/golang/mock/gomock
//go get -u github.com/golang/mock/mockgen

//第一步：使用 mockgen 生成 db_mock.go。一般传递三个参数。包含需要被mock的接口得到源文件source，生成的目标文件destination，包名package。

//$ mockgen -source=db.go -destination=db_mock.go -package=main
//在上面的例子中，当 Get() 的参数为 Tom，则返回 error，这称之为打桩(stub)，有明确的参数和返回值是最简单打桩方式。除此之外，检测调用次数、调用顺序，动态设置返回值等方式也经常使用。

//3.1 参数(Eq, Any, Not, Nil)

//m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
//m.EXPECT().Get(gomock.Any()).Return(630, nil)
//m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
//m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil"))
//Eq(value) 表示与 value 等价的值。
//Any() 可以用来表示任意的入参。
//Not(value) 用来表示非 value 以外的值。
//Nil() 表示 None 值
//3.2 返回值(Return, DoAndReturn)

//m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
//m.EXPECT().Get(gomock.Any()).Do(func(key string) {
//t.Log(key)
//})
//m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
//if key == "Sam" {
//return 630, nil
//}
//return 0, errors.New("not exist")
//})
//Return 返回确定的值
//Do Mock 方法被调用时，要执行的操作吗，忽略返回值。
//DoAndReturn 可以动态地控制返回值。
//3.3 调用次数(Times)

//func TestGetFromDB(t *testing.T) {
//ctrl := gomock.NewController(t)
//defer ctrl.Finish()

//m := NewMockDB(ctrl)
//m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
//GetFromDB(m, "ABC")
//GetFromDB(m, "DEF")
//}
//Times() 断言 Mock 方法被调用的次数。
//MaxTimes() 最大次数。
//MinTimes() 最小次数。
//AnyTimes() 任意次数（包括 0 次）。
//3.4 调用顺序(InOrder)
//func TestGetFromDB(t *testing.T) {
//ctrl := gomock.NewController(t)
//defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

//m := NewMockDB(ctrl)
//o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
//o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
//gomock.InOrder(o1, o2)
//GetFromDB(m, "Tom")
//GetFromDB(m, "Sam")
//}
//4 如何编写可 mock 的代码
//写可测试的代码与写好测试用例是同等重要的，如何写可 mock 的代码呢？

//mock 作用的是接口，因此将依赖抽象为接口，而不是直接依赖具体的类。
//不直接依赖的实例，而是使用依赖注入降低耦合性。
