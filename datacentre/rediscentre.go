package datacentre

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	cmc "github.com/zytzjx/anthenacmc/cmcserverinfo"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

// ExampleClient : testgu
// Printer : what is this?
// Greet : describe what this function does
// CreateMessage : describe what this function does
func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	err = rdb.HSet(ctx, "runoobkey", "JJJ", "aassaa").Err()
	if err != nil {
		panic(err)
	}

	dict, err := rdb.HGetAll(ctx, "runoobkey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("runoobkey", dict)

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// user := map[string]string{"Name": "Pradeep", "Company": "SCTL", "Address": "Mumbai", "Location": "RCP"}
	RedisPUBSUB(rdb)

	// Output: key value
	// key2 does not exist
}

//RedisPUBSUB PUBSUB
func RedisPUBSUB(rdb *redis.Client) {
	pubsub := rdb.PSubscribe(ctx, "mychannel*")
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// // Publish a message.
	// err = rdb.Publish(ctx, "mychannel1ABC", "hello").Err()
	// if err != nil {
	// 	panic(err)
	// }

	// time.AfterFunc(time.Second, func() {
	// 	// When pubsub is closed channel is closed too.
	// 	_ = pubsub.Close()
	// })

	// Consume messages.
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}

// SaveSerialConfigRedis Save SerialConfig to Redis DB
func SaveSerialConfigRedis(confInstall cmc.ConfigInstall) {
	file, _ := json.Marshal(confInstall.Results[0])
	var result map[string]string
	if err := json.Unmarshal(file, &result); err != nil {
		fmt.Println("data error")
	}

	for k, v := range result { // HGETALL serialconfig
		rdb.HMSet(ctx, "serialconfig", k, v)
	}
	// rdb.HGetAll(ctx, "serialconfig").Result()
}

// GetSerialConfig getserial config
func GetSerialConfig() (cmc.ConfigResult, error) {
	var result cmc.ConfigResult
	cr, err := rdb.HGetAll(ctx, "serialconfig").Result()
	if err != nil {
		return result, err
	}
	dd, _ := json.Marshal(cr)

	if err := json.Unmarshal(dd, &result); err != nil {
		fmt.Println("data error")
		return result, err
	}
	return result, nil
}
