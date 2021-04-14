package main

import (
	"encoding/json"
	"fmt"
	"github.com/Nubes3/common"
	"github.com/Nubes3/common/models/arangodb"
	"github.com/Nubes3/common/models/nats"
	"github.com/gin-gonic/gin"
	n "github.com/nats-io/nats.go"
	"net/http"
	"time"
)

func main() {
	common.InitCoreComponents(false, false, false, true)
	r := gin.New()

	sub1, _ := nats.Nc.QueueSubscribe(nats.UserSubj, "userQueue", func(msg *n.Msg) {
		fmt.Println(string(msg.Data))

		u := arangodb.UserRes{
			Firstname: "Test",
			Lastname:  "Test",
			Username:  "12345",
			Pass:      "1232412421321",
			Email:     "132141254124",
			Dob:       time.Now(),
			Company:   "Test C",
			Gender:    false,
		}

		userJson, _ := json.Marshal(u)
		res := nats.MsgResponse{
			IsErr:     false,
			Data:      string(userJson),
			ExtraData: []string{"123", "111"},
		}

		resJson, _ := json.Marshal(res)
		_ = msg.Respond(resJson)
	})

	sub2, _ := nats.Nc.QueueSubscribe(nats.UserSubj, "userQueue", func(msg *n.Msg) {
		fmt.Println(string(msg.Data))

		u := arangodb.UserRes{
			Firstname: "Test 1",
			Lastname:  "Test 1",
			Username:  "123456",
			Pass:      "12324124213211111",
			Email:     "132141254121114",
			Dob:       time.Now(),
			Company:   "Test C1",
			Gender:    false,
		}

		userJson, _ := json.Marshal(u)
		res := nats.MsgResponse{
			IsErr:     false,
			Data:      string(userJson),
			ExtraData: []string{"1213", "1111"},
		}

		resJson, _ := json.Marshal(res)
		_ = msg.Respond(resJson)
	})

	defer sub1.Unsubscribe()
	defer sub2.Unsubscribe()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "kay")
	})

	r.GET("/testNats", func(c *gin.Context) {
		msg := nats.Msg{
			ReqType:   nats.Resolve,
			Data:      "1e12321412ad",
			ExtraData: nil,
		}
		jsonData, _ := json.Marshal(msg)
		repRaw, _ := nats.Nc.Request(nats.UserSubj, jsonData, time.Second*10)

		var reply nats.MsgResponse
		_ = json.Unmarshal(repRaw.Data, &reply)
		var receivedUser arangodb.User
		_ = json.Unmarshal([]byte(reply.Data), &receivedUser)

		c.JSON(http.StatusOK, receivedUser)
	})

	r.Run(":7179")
}
