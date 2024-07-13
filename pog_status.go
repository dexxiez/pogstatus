package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	wifiSignalChan := make(chan int)

	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer X.Close()

	go updateWifiSignal(wifiSignalChan)

	wifiSignal := -1

	for {
		select {
		case newSignal := <-wifiSignalChan:
			wifiSignal = newSignal
		default:
		}

		status := fmt.Sprintf("CPU: %d%% | Mem: %d%% | Wifi: %d%% | %s", getCPUUsage(), getMemUsage(), wifiSignal, time.Now().Format("03:04:05 PM"))
		setRootName(X, status)
		time.Sleep(1 * time.Second)
	}
}

func setRootName(X *xgb.Conn, name string) {
	root := xproto.Setup(X).DefaultScreen(X).Root
	xproto.ChangeProperty(X, xproto.PropModeReplace, root, xproto.AtomWmName,
		xproto.AtomString, 8, uint32(len(name)), []byte(name))
}

func getCPUUsage() int {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("CPU error!")
		log.Fatal(err)
	}
	return int(percent[0])
}

func getMemUsage() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Memory error!")
		log.Fatal(err)
	}

	return int(v.UsedPercent)
}

func updateWifiSignal(signalChan chan<- int) {
	for {
		cmd := exec.Command("sh", "-c", "nmcli -t -f IN-USE,SIGNAL dev wifi | grep '*' | awk -F: '{print $2}'")
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Error executing command: %v", err)
			continue
		}

		signalStr := strings.TrimSpace(string(output))
		signal, err := strconv.Atoi(signalStr)
		if err != nil {
			log.Printf("Error converting signal strength to int: %v", err)
			continue
		}

		signalChan <- signal
		time.Sleep(10 * time.Second)
	}
}
