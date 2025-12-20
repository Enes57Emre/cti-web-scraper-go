package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func fetchHTML(url, outputFile string) error {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("istek oluşturulamadı: %w", err)
	}

	req.Header.Set("User-Agent", "CTI-Scraper/1.0 (+https://example.com)")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("istek gönderilirken hata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("beklenmeyen HTTP durum kodu: %d %s", resp.StatusCode, resp.Status)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("çıktı dosyası oluşturulamadı: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("HTML içerik dosyaya yazılamadı: %w", err)
	}

	return nil
}

func takeScreenshot(url, screenshotFile string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var imgBuf []byte

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(3 * time.Second),
		chromedp.FullScreenshot(&imgBuf, 90),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return fmt.Errorf("ekran görüntüsü alınırken hata: %w", err)
	}

	if err := os.WriteFile(screenshotFile, imgBuf, 0644); err != nil {
		return fmt.Errorf("ekran görüntüsü dosyaya yazılamadı: %w", err)
	}

	return nil
}

func extractLinksFromHTML(htmlFile, urlsFile string) error {
	f, err := os.Open(htmlFile)
	if err != nil {
		return fmt.Errorf("HTML dosyası açılamadı: %w", err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return fmt.Errorf("HTML parse edilirken hata: %w", err)
	}

	out, err := os.Create(urlsFile)
	if err != nil {
		return fmt.Errorf("URL çıktısı için dosya oluşturulamadı: %w", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	seen := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		if seen[href] {
			return
		}
		seen[href] = true

		_, _ = writer.WriteString(href + "\n")
	})

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go <hedef_URL>")
		os.Exit(1)
	}

	url := os.Args[1]
	htmlOutput := "output.html"
	screenshotOutput := "screenshot.png"
	urlsOutput := "urls.txt"

	log.Printf("Hedef URL: %s\n", url)

	log.Println("HTML içeriği indiriliyor...")
	if err := fetchHTML(url, htmlOutput); err != nil {
		log.Fatalf("HTML çekme hatası: %v\n", err)
	}
	log.Printf("HTML içerik '%s' dosyasına kaydedildi.\n", htmlOutput)

	log.Println("Ekran görüntüsü alınıyor...")
	if err := takeScreenshot(url, screenshotOutput); err != nil {
		log.Fatalf("Ekran görüntüsü hatası: %v\n", err)
	}
	log.Printf("Ekran görüntüsü '%s' dosyasına kaydedildi.\n", screenshotOutput)

	log.Println("Sayfadaki linkler çıkarılıyor...")
	if err := extractLinksFromHTML(htmlOutput, urlsOutput); err != nil {
		log.Fatalf("Link çıkarma hatası: %v\n", err)
	}
	log.Printf("Bulunan linkler '%s' dosyasına kaydedildi.\n", urlsOutput)

	log.Println("İşlem başarıyla tamamlandı.")
}
