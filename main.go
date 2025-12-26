package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"tor-scraper/internal/config"
	"tor-scraper/internal/output"
	"tor-scraper/internal/proxy"
	"tor-scraper/internal/scraper"
)

func main() {
	configPath := flag.String("config", "targets.yaml", "Hedef listesi YAML dosyasının yolu")
	checkTor := flag.Bool("check", false, "Sadece Tor bağlantısını kontrol et")
	flag.Parse()

	printBanner()

	fmt.Printf("[INFO] Yapılandırma dosyası yükleniyor: %s\n", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("[ERROR] Yapılandırma yüklenemedi: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[INFO] %d hedef bulundu\n", len(cfg.Targets))

	fmt.Printf("[INFO] Tor proxy'ye bağlanılıyor (%s:%s)...\n", cfg.Proxy.Host, cfg.Proxy.Port)
	torClient, err := proxy.NewTorClient(cfg.Proxy.Host, cfg.Proxy.Port, cfg.Timeout)
	if err != nil {
		fmt.Printf("[ERROR] Tor client oluşturulamadı: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[INFO] Tor bağlantısı kontrol ediliyor...")
	isConnected, msg, err := torClient.CheckTorConnection()
	if err != nil {
		fmt.Printf("[WARN] Tor bağlantısı doğrulanamadı: %v\n", err)
		fmt.Println("[WARN] Tarama yine de devam edecek...")
	} else if isConnected {
		fmt.Printf("[INFO] %s\n", msg)
	}

	if *checkTor {
		if isConnected {
			fmt.Println("[SUCCESS] Tor bağlantısı başarılı!")
			os.Exit(0)
		} else {
			fmt.Println("[FAILED] Tor bağlantısı kurulamadı!")
			os.Exit(1)
		}
	}

	writer, err := output.NewWriter(".")
	if err != nil {
		fmt.Printf("[ERROR] Output writer oluşturulamadı: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	fmt.Printf("[INFO] Çıktılar kaydedilecek: %s\n", writer.GetOutputDir())

	s := scraper.NewScraper(torClient, writer)
	startTime := time.Now()

	s.ScanTargets(cfg.Targets)

	results := s.GetResults()
	total, success, failed, timeout := s.GetSummary()

	var entries []output.ScanReportEntry
	for _, r := range results {
		entry := output.ScanReportEntry{
			Name:      r.Name,
			URL:       r.URL,
			Status:    r.Status,
			Message:   r.Message,
			Timestamp: r.Timestamp,
			Duration:  r.Duration.String(),
		}
		entries = append(entries, entry)
		writer.WriteLog(r.Status, r.Name, r.URL, r.Message, r.Duration)
	}

	err = writer.SaveJSONReport(entries)
	if err != nil {
		fmt.Printf("[WARN] JSON raporu kaydedilemedi: %v\n", err)
	}

	writer.WriteSummary(total, success, failed, timeout)

	fmt.Printf("\n[INFO] Toplam süre: %v\n", time.Since(startTime).Round(time.Second))
	fmt.Printf("[INFO] Raporlar: %s\n", writer.GetOutputDir())
}

func printBanner() {
	banner := `
============================================================
                TOR SCRAPER by Healme                                                    
============================================================`
	fmt.Println(banner)
}
