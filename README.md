# TOR Scraper

Dark Web istihbarat toplama aracı. Bu araç, .onion adreslerini Tor ağı üzerinden tarayarak veri toplar.

## Özellikler

- **YAML Tabanlı Hedef Listesi**: Hedefler kolayca yapılandırılabilir
- **Tor SOCKS5 Proxy Desteği**: Tüm trafik Tor ağı üzerinden yönlendirilir
- **Otomatik Hata Yönetimi**: Çalışmayan siteler programı durdurmaz
- **Detaylı Raporlama**: JSON ve log formatında raporlar
- **Tarih Damgalı Çıktılar**: Her tarama ayrı bir klasörde saklanır

## Kurulum

### Gereksinimler

1. **Go 1.21+** kurulu olmalı
2. **Tor servisi** arka planda çalışıyor olmalı

### Tor Servisini Başlatma

```bash
# Linux (apt)
sudo apt install tor
sudo systemctl start tor

# macOS
brew install tor
brew services start tor

# Windows
# Tor Browser'ı açın (otomatik olarak 9150 portunda çalışacak)
```

### Projeyi Derleme

```bash
go mod tidy

# Projeyi derle
go build -o tor-scraper main.go

# veya doğrudan çalıştır
go run main.go
```

## Kullanım

### Temel Kullanım

```bash
# Varsayılan
go run main.go

# Özel yaml dosyası kullanılabilir
go run main.go -config hedefler.yaml

# Tor bağlantısını kontrol et
go run main.go -check
```

### Komut Satırı Argümanları

| Argüman | Açıklama | Varsayılan |
|---------|----------|------------|
| `-config` | Hedef listesi YAML dosyası | `targets.yaml` |
| `-check` | Sadece Tor bağlantısını kontrol et | `false` |

## Yapılandırma

### targets.yaml Örneği

```yaml
# Proxy
proxy:
  host: "127.0.0.1"
  port: "9050"  # Tor Browser için 9150

timeout: 60

targets:
  - name: "X"
    url: "http://Xdomain.onion"
  
  - name: "Y"
    url: "http://Ydomain.onion"
```

## Çıktılar

Her tarama sonucunda aşağıdaki çıktılar oluşturulur:

```
output/
└── 2024-01-15_14-30-00/
    ├── html/
    │   ├── Y.html
    │   └── X.html
    ├── scan_report.log
    └── scan_results.json
```

### scan_report.log
```
=== TOR SCRAPER TARAMA RAPORU ===
Tarih: 2024-01-15 14:30:00

[SUCCESS] 16:38:20 | Not Evil | 4794 byte veri alındı | 2.901s

[ERROR] 16:38:20 | Mazafaka | HTTP Durum Kodu: 403 | 4.186s

========================================
TARAMA ÖZETİ
========================================
Toplam Hedef    : 2
Başarılı        : 1
Başarısız       : 1
Zaman Aşımı     : 0
Başarı Oranı    : 50.0%
========================================
```

### scan_results.json
```json
[
  {
    "name": "Alte***",
    "url": "http://xxxxx.onion/index.php",
    "status": "SUCCESS",
    "message": "2399 byte veri alındı",
    "timestamp": "2025-12-26T16:37:55.813720959+03:00",
    "duration": "3.885663125s"
  }
]
```

## Proje Yapısı

```
tor-scraper/
├── main.go                 # Ana uygulama
├── go.mod                  # Go modül dosyası
├── targets.yaml            # Hedef listesi
├── README.md               # Bu dosya
└── internal/
    ├── config/
    │   └── config.go       # YAML yapılandırma okuyucu
    ├── proxy/
    │   └── proxy.go        # Tor SOCKS5 client
    ├── scraper/
    │   └── scraper.go      # Tarama mantığı
    └── output/
        └── output.go       # Çıktı yazıcı
```

## IP Sızıntısı Kontrolü

Program çalışmadan önce Tor bağlantısını otomatik olarak kontrol eder. Manuel kontrol için:

```bash
go run main.go -check
```

## 

*Bu proje eğitim amaçlıdır*

