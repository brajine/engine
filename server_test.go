package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"mtlive/data"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

func run() {
	time.Sleep(500 * time.Millisecond)

	ord1 := data.OrderType{
		Symbol:    []byte("eurusd"),
		TimeOpen:  []byte("2019.10.16 14:50:42"),
		Type:      []byte("OP_SELL"),
		CurVolume: []byte("0.1"),
		PriceOpen: []byte("1.10266"),
		Swap:      []byte("-0.75"),
		Profit:    []byte("-15.02"),
	}

	ord2 := data.OrderType{
		Symbol:    []byte("gbpusd"),
		TimeOpen:  []byte("2019.10.16 34:34:34"),
		Type:      []byte("OP_BUYSTOP"),
		CurVolume: []byte("0.2"),
		PriceOpen: []byte("4.10266"),
	}

	msg := data.TradesMsg{}

	msg.Server = []byte("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	msg.Company = []byte("11111")

	msg.Orders = make(map[string]data.OrderType)
	msg.Orders["466033934"] = ord1
	msg.Orders["466033935"] = ord2

	conn, err := net.Dial("tcp", "127.0.0.1:3131")
	if err == nil {
		enc := gob.NewEncoder(conn)
		enc.Encode(msg)
	}
	conn.Close()
}

func TestMain(t *testing.T) {
	go run()

	ln, err := net.Listen("tcp", ":3131")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	// Read message
	var f int
	for {
		if conn, err := ln.Accept(); err == nil {
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if n > 0 {
					ioutil.WriteFile("data_"+strconv.Itoa(f)+".bin", buf[:n], os.ModePerm)
					f++
					for i := 0; i < n; i++ {
						fmt.Printf("%02X ", buf[i])
					}
					fmt.Println()
					fmt.Println()
				}
				if err != nil {
					break
				}
			}

			// conn.Close()
		}

		// Send message
		// upd := Response{
		// 	Msg: 1,
		// }

		// for {
		// 	if conn, err := ln.Accept(); err == nil {
		// 		enc := gob.NewEncoder(conn)
		// 		enc.Encode(upd)
		// 		conn.Close()
		// 	}
		// }
	}

}
