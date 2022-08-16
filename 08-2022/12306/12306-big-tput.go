
...
//localSpike包结构体定义
package localSpike

type LocalSpike struct {
    LocalInStock     int64
    LocalSalesVolume int64
}
...
//remoteSpike对hash结构的定义和redis连接池
package remoteSpike
//远程订单存储健值
type RemoteSpikeKeys struct {
    SpikeOrderHashKey string    //redis中秒杀订单hash结构key
    TotalInventoryKey string    //hash结构中总订单库存key
    QuantityOfOrderKey string   //hash结构中已有订单数量key
}

//初始化redis连接池
func NewPool() *redis.Pool {
    return &redis.Pool{
        MaxIdle:   10000,
        MaxActive: 12000, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", ":6379")
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    }
}
...
func init() {
    localSpike = localSpike2.LocalSpike{
        LocalInStock:     150,
        LocalSalesVolume: 0,
    }
    remoteSpike = remoteSpike2.RemoteSpikeKeys{
        SpikeOrderHashKey:  "ticket_hash_key",
        TotalInventoryKey:  "ticket_total_nums",
        QuantityOfOrderKey: "ticket_sold_nums",
    }
    redisPool = remoteSpike2.NewPool()
    done = make(chan int, 1)
    done <- 1
}
