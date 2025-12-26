package proxy

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type TorClient struct {
	httpClient *http.Client
	proxyAddr  string
}

func NewTorClient(host, port string, timeout int) (*TorClient, error) {
	proxyAddr := fmt.Sprintf("%s:%s", host, port)

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("SOCKS5 proxy oluşturulamadı: %w", err)
	}

	transport := &http.Transport{
		Dial:                  dialer.Dial,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	return &TorClient{
		httpClient: httpClient,
		proxyAddr:  proxyAddr,
	}, nil
}

func (tc *TorClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0")

	return tc.httpClient.Do(req)
}

func (tc *TorClient) CheckTorConnection() (bool, string, error) {
	resp, err := tc.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "Tor bağlantısı doğrulandı", nil
	}

	return false, "Tor bağlantısı doğrulanamadı", nil
}

func (tc *TorClient) GetProxyAddr() string {
	return tc.proxyAddr
}
