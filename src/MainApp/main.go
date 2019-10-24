package main

import (
	"modbuslib"
	"fmt"
	"time"
)

func main() {
	client := new(modbuslib.ModbusClient)
	client.Host = "192.168.11.50:502"

	_, err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	go func(client *modbuslib.ModbusClient, period int64) {
		// period는 Millisecond 단위
		var t time.Time = time.Now()

		for {
			data, err := client.ReadHoldingRegister(3012, 6)
			if err != nil {
				fmt.Println(err)
				return
			}
			
			
			fmt.Println(time.Now(), data)
			time.Sleep(time.Duration(period - time.Since(t).Nanoseconds() / 1000000 % period) * time.Millisecond)
		}
	}(client, 100)

	fmt.Scanln()
}