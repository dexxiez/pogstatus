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
		case wifiSignal = <-wifiSignalChan:
		default:
		}

		status := fmt.Sprintf(" CPU: %s%% | Mem: %s%% | Wifi: %s | %s ",
			padLengthSpaces(getCPUUsage(), 2),
			padLengthSpaces(getMemUsage(), 2),
			getWifiSignalIcon(wifiSignal),
			getDateTime(),
		)
		setRootName(X, status)
		time.Sleep(time.Second)
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

func getDateTime() string {
	return time.Now().Format("02/01/2006 03:04:05 PM")
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

func padLengthSpaces(n int, length int) string {
	s := strconv.Itoa(n)
	spaces := length - len(s)
	if spaces > 0 {
		return strings.Repeat(" ", spaces) + s
	}
	return s
}

func getWifiSignalIcon(signal int) string {
	if signal < 0 {
		return " "
	} else if signal < 20 {
		return " "
	} else if signal < 40 {
		return " "
	} else if signal < 60 {
		return " "
	} else if signal < 80 {
		return " "
	} else {
		return " "
	}
}
