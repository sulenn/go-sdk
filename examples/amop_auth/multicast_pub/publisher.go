package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"
)

func onPush(data []byte) {
	log.Printf("received: %s\n", string(data))
}

const (
	publicKey1 = "7cd0596006e3c0549482d010735303f25d6c308ef66643b63deef0a382a7620556eb49641cb3d45f4901a068a5e68475f8ba3b03a1bc785e84fe6490d66df365"
	publicKey2 = "19dece101df106ca4baf478f98911cdc525db5c6b58f2189af9f69ff314e9f0bcb816b41fb8bd49ae830dc1087bf51c71a21c3e3a332132262b5ecf0189817f4"
	publicKey3 = "8b38138ea887220289276ca700e162647af79b4c61f33aefcfdaa2c3b714b2983084e519273208e8646b7f840e91b9053952df28a3bce1a6bca0132c26a36694"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("the number of arguments is not equal 4")
	}
	endpoint := os.Args[1]
	topic := os.Args[2]
	config := &conf.Config{IsHTTP: false, ChainID: 1, CAFile: "ca.crt", Key: "sdk.key", Cert: "sdk.crt", IsSMCrypto: false, GroupID: 1,
		PrivateKey: "145e247e170ba3afd6ae97e88f00dbc976c2345d511b0f6713355d19d8b80b58",
		NodeURL:    endpoint}
	c, err := client.Dial(config)
	if err != nil {
		log.Fatalf("init publisher failed, err: %v\n", err)
	}
	publicKeys := make([]*ecdsa.PublicKey, 0)
	pubKey1, err := crypto.UnmarshalPubkey(hexutils.HexToBytes("04" + publicKey1))
	pubKey2, err := crypto.UnmarshalPubkey(hexutils.HexToBytes("04" + publicKey2))
	pubKey3, err := crypto.UnmarshalPubkey(hexutils.HexToBytes("04" + publicKey3))
	if err != nil {
		log.Fatalf("decompress pubkey failed, err: %v", err)
	}
	publicKeys = append(publicKeys, pubKey1, pubKey2, pubKey3)
	err = c.PublishPrivateTopic(topic, publicKeys, onPush)
	if err != nil {
		log.Fatalf("publish topic failed, err: %v\n", err)
	}
	fmt.Println("publish topic success")
	time.Sleep(3 * time.Second)

	message := "Hi, FISCO BCOS!"
	go func() {
		for i := 0; i < 1000; i++ {
			log.Printf("publish message: %s ", message+" "+strconv.Itoa(i))
			err = c.BroadcastAMOPPrivateMsg(topic, []byte(message+" "+strconv.Itoa(i)))
			time.Sleep(2 * time.Second)
			if err != nil {
				log.Printf("PushTopicDataRandom failed, err: %v\n", err)
			}
		}
	}()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
}
