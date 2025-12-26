package scraper

import (
	"fmt"
	"io"
	"time"

	"tor-scraper/internal/config"
	"tor-scraper/internal/output"
	"tor-scraper/internal/proxy"
)

type ScanResult struct {
	Name      string
	URL       string
	Status    string
	Message   string
	Timestamp time.Time
	Duration  time.Duration
}

type Scraper struct {
	client  *proxy.TorClient
	writer  *output.Writer
	results []ScanResult
}

func NewScraper(client *proxy.TorClient, writer *output.Writer) *Scraper {
	return &Scraper{
		client:  client,
		writer:  writer,
		results: make([]ScanResult, 0),
	}
}

func (s *Scraper) ScanTargets(targets []config.Target) {
	fmt.Printf("\n[INFO] Toplam %d hedef taranacak\n", len(targets))
	fmt.Println("========================================")

	for i, target := range targets {
		fmt.Printf("\n[%d/%d] Taranıyor: %s\n", i+1, len(targets), target.URL)
		result := s.scanSingleTarget(target)
		s.results = append(s.results, result)

		s.printResult(result)
	}

	fmt.Println("\n========================================")
	fmt.Println("[INFO] Tarama tamamlandı!")
}

func (s *Scraper) scanSingleTarget(target config.Target) ScanResult {
	startTime := time.Now()

	result := ScanResult{
		Name:      target.Name,
		URL:       target.URL,
		Timestamp: startTime,
	}

	resp, err := s.client.Get(target.URL)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Status = "ERROR"
		result.Message = err.Error()

		if isTimeout(err) {
			result.Status = "TIMEOUT"
			result.Message = "Bağlantı zaman aşımına uğradı"
		}

		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("HTTP Durum Kodu: %d", resp.StatusCode)
		return result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("İçerik okunamadı: %v", err)
		return result
	}

	err = s.writer.SaveHTML(target.Name, target.URL, body)
	if err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("Dosya kaydedilemedi: %v", err)
		return result
	}

	result.Status = "SUCCESS"
	result.Message = fmt.Sprintf("%d byte veri alındı", len(body))

	return result
}

func (s *Scraper) printResult(result ScanResult) {
	fmt.Printf("[%s] %s -> %s (%v)\n", result.Status, result.Name, result.Message, result.Duration.Round(time.Millisecond))
}

func (s *Scraper) GetResults() []ScanResult {
	return s.results
}

func (s *Scraper) GetSummary() (total, success, failed, timeout int) {
	total = len(s.results)
	for _, r := range s.results {
		switch r.Status {
		case "SUCCESS":
			success++
		case "TIMEOUT":
			timeout++
		default:
			failed++
		}
	}
	return
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "timeout") || contains(errStr, "deadline exceeded")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
