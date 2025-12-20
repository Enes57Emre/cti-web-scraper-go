# CTI Web Scraper (Go) – HTML + Screenshot + URL List

Bu proje, Siber Tehdit İstihbaratı (CTI) temel gereksinimi olarak tek bir web sayfasının:
- Ham HTML içeriğini indirmeyi ve dosyaya kaydetmeyi,
- Sayfanın ekran görüntüsünü almayı,
- (Ek görev) sayfadaki bağlantıları (a[href]) listelemeyi
amaçlayan basit bir Go (Golang) aracıdır.

## Özellikler
- Komut satırından URL alır
- HTTP GET ile HTML içeriğini kaydeder
- chromedp ile tam ekran ekran görüntüsü üretir
- goquery ile sayfadaki linkleri `urls.txt` dosyasına yazar

## Gereksinimler
- Go (Golang)
- Google Chrome / Chromium (chromedp için)

## Kurulum
Proje klasöründe:

```bash
go mod init cti-scraper
go get github.com/chromedp/chromedp
go get github.com/PuerkitoBio/goquery
