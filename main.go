package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	gw "mqttsngws/gateway"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

const Version string = "7.0.0.1"

func main() {
	var (
		confFile  = flag.String("c", "", "config file path")
		topicFile = flag.String("t", "", "predefined topic file path")
	)
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if dir == "/" {
		dir += "data"
	} else {
		dir += "/data"
	}
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	if *confFile == "" {
		*confFile = dir + "/gateway.yml"
	}

	// parse config
	config, err := gw.ParseConfig(*confFile)
	if err != nil {
		panic(err)
	}

	// parse topic file
	if *topicFile != "" {
		err = gw.InitPredefinedTopic(*topicFile)
		if err != nil {
			panic(err)
		}
	}
	log_file_path := dir + "/" + config.LogFilePath
	// initialize logger
	err = gw.InitLogger(log_file_path)
	if err != nil {
		panic(err)
	}

	go func() {
		var uip string
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if strings.HasPrefix(ipnet.IP.String(), "192.168.") {
					if ipnet.IP.To4() != nil {
						//dns = ipnet.IP.To4()
						ips := strings.Split(ipnet.IP.String(), ".")
						ips[3] = "255"
						uip = strings.Join(ips, ".")
					}
				}
			}
		}
		ip := net.ParseIP(uip)
		srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
		dstAddr := &net.UDPAddr{IP: ip, Port: 13623}
		umsg := make(map[string]interface{})
		umsg["v"] = Version
		umsg["tcp"] = config.BrokerPort
		umsg["tls"] = config.BrokerPort
		umsg["ag"] = config.Port
		pmsg, err := json.Marshal(umsg)
		if err != nil {
			fmt.Println(err)
		}
		ticker := time.NewTicker(10 * time.Second)
		go func() {
			for range ticker.C {
				conn, err := net.ListenUDP("udp", srcAddr)
				if err != nil {
					fmt.Println(err)
					conn.Close()
				}
				//send descover package
				conn.WriteToUDP([]byte(pmsg), dstAddr)
			}
		}()
	}()

	// create signal chan
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// create Gateway
	var gateway gw.Gateway
	if config.IsAggregate {
		gateway = gw.NewAggregatingGateway(config, signalChan)
	} else {
		gateway = gw.NewTransparentGateway(config, signalChan)
	}

	// start server
	err = gateway.StartUp()
	if err != nil {
		log.Println(errors.New("ERROR : failed to StartUp gateway"))
	}
}
