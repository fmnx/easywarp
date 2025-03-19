package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/fmnx/easywarp/tunsetup"
	"github.com/tidwall/gjson"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Warp struct {
	Auto       bool   `yaml:"auto" json:"auto"`
	Endpoint   string `yaml:"endpoint" json:"endpoint"`
	IPv4       string `yaml:"ipv4" json:"ipv4"`
	IPv6       string `yaml:"ipv6" json:"ipv6"`
	PrivateKey string `yaml:"private-key" json:"private-key"`
	PublicKey  string `yaml:"public-key" json:"public-key"`
	ClientID   []byte `yaml:"client-id" json:"client-id"`
}

func (w *Warp) load() {
	buf, err := os.ReadFile(".warp.json")
	if err != nil {
		w.apply()
		w.save()
	} else {
		_ = json.Unmarshal(buf, w)
	}

}

func (w *Warp) save() {
	warpFile, _ := json.MarshalIndent(w, "", "  ")
	_ = os.WriteFile(".warp.json", warpFile, 0644)
}

func (w *Warp) apply() {

	log.Println("Automatically applying for Warp...")

	url := "https://api.cloudflareclient.com/v0a2223/reg"

	privateKey := NewPrivateKey()
	// 请求头
	headers := map[string]string{
		"CF-Client-Version": "a-6.11-2223",
		"Host":              "api.cloudflareclient.com",
		"Connection":        "Keep-Alive",
		"Accept-Encoding":   "gzip",
		"User-Agent":        "okhttp/3.12.1",
		"Content-Type":      "application/json",
	}

	jsonData, _ := json.Marshal(map[string]string{
		"key":    privateKey.Public().String(),
		"locale": "en-US",
		"tos":    time.Now().Format(time.RFC3339Nano),
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("A request error occurred while automatically applying for Warp.: %v\n", err)
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	body, _ := io.ReadAll(reader)
	clientID := gjson.Get(string(body), "config.client_id").String()
	ipv6 := gjson.Get(string(body), "config.interface.addresses.v6").String()
	if ipv6 == "" {
		log.Fatalln("Failed to automatically apply for Warp.")
	}

	w.ClientID, _ = base64.StdEncoding.DecodeString(clientID)
	w.IPv4 = "172.16.0.2"
	w.IPv6 = ipv6
	w.PrivateKey = privateKey.String()
	w.PublicKey = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
	w.Endpoint = "engage.cloudflareclient.com:2408"

	log.Println("Warp has been successfully applied.")
}

func (w *Warp) Run() {
	w.load()

	tunDev, err := tun.CreateTUN("wg0", 1280)
	if err != nil {
		log.Fatalf("Failed to create TUN device: wg0, Error: %s\n", err.Error())
	}
	err = tunsetup.ConfigureTunAddr("wg0", w.IPv4, w.IPv6)
	if err != nil {
		log.Fatalln(err.Error())
	}

	bind := conn.NewDefaultBind()
	bind.SetReserved(w.ClientID)
	logger := device.NewLogger(1, "")
	dev := device.NewDevice(tunDev, bind, logger)
	err = dev.Up()
	if err != nil {
		log.Fatalf("Failed to bring up device: %v\n", err)
	}
	dev.SetPrivateKey(w.PrivateKey)
	peer := dev.SetPublicKey(w.PublicKey)

	dev.SetEndpoint(peer, resolvEndpoint(w.Endpoint)).SetAllowedIP(peer)
	peer.HandlePostConfig()
}

func resolvEndpoint(endpoint string) string {
	c, err := net.DialTimeout("udp", endpoint, 3*time.Second)
	defer c.Close()
	if err != nil {
		return "162.159.192.1:2408"
	} else {
	}
	return c.RemoteAddr().String()
}
