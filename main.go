package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	// یہ بلینک امپورٹ گو کو بتائے گا کہ whatsmeow کو ڈیلیٹ نہیں کرنا!
	_ "go.mau.fi/whatsmeow"
)

func generateZipDocs() error {
	zipFileName := "whatsmeow_docs.zip"
	
	// زپ فائل کریئیٹ کریں
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("زپ فائل بنانے میں مسئلہ: %v", err)
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	fmt.Println("\n🚀 ----------------------------------------------------")
	fmt.Println("🚀 STEP 1: واٹس میو کے تمام پیکجز تلاش کیے جا رہے ہیں...")
	fmt.Println("🚀 ----------------------------------------------------")

	cmdList := exec.Command("go", "list", "go.mau.fi/whatsmeow/...")
	var listErr bytes.Buffer
	cmdList.Stderr = &listErr
	outputList, err := cmdList.Output()

	if err != nil {
		fmt.Printf("❌ ایرر 'go list' کمانڈ چلانے میں: %v\n", err)
		return fmt.Errorf("پیکجز فائنڈ نہیں ہوئے")
	}

	packagesStr := strings.TrimSpace(string(outputList))
	packages := strings.Split(packagesStr, "\n")
	fmt.Printf("✅ ٹوٹل %d پیکجز مل گئے۔ زپ فائل بنائی جا رہی ہے...\n\n", len(packages))

	for i, pkg := range packages {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			continue
		}

		fmt.Printf("⏳ [%d/%d] پروسیسنگ ہو رہی ہے: %s\n", i+1, len(packages), pkg)

		cmdDoc := exec.Command("go", "doc", "-all", pkg)
		outputDoc, err := cmdDoc.Output()

		if err != nil || len(outputDoc) == 0 {
			fmt.Printf("  ⚠️ کوئی پبلک فنکشن نہیں ملا، اسے چھوڑ دیا گیا۔\n")
			continue
		}

		// فائل کا نام بنانا (مثلاً: go.mau.fi/whatsmeow/types کو whatsmeow_types.txt بنا دے گا)
		fileName := strings.TrimPrefix(pkg, "go.mau.fi/")
		fileName = strings.ReplaceAll(fileName, "/", "_") + ".txt"

		// زپ کے اندر نئی ٹیکسٹ فائل بنانا
		f, err := zipWriter.Create(fileName)
		if err != nil {
			fmt.Printf("  ❌ زپ کے اندر فائل بنانے میں ایرر: %v\n", err)
			continue
		}

		// اس ٹیکسٹ فائل میں سارا ڈیٹا لکھنا
		_, err = f.Write(outputDoc)
		if err != nil {
			fmt.Printf("  ❌ ڈیٹا لکھنے میں ایرر: %v\n", err)
		} else {
			fmt.Printf("  ✅ کامیابی! %s زپ میں شامل کر دی گئی۔\n", fileName)
		}
	}

	fmt.Println("\n🎉 سب کچھ مکمل ہو گیا! زپ فائل ریڈی ہے۔")
	return nil
}

func main() {
	fmt.Println("=================================================")
	fmt.Println("   WHATSMEOW ZIP EXTRACTOR SERVER STARTING...    ")
	fmt.Println("=================================================")

	if err := generateZipDocs(); err != nil {
		fmt.Printf("\n🚨 خطرہ: ڈاکیومنٹیشن جنریٹ کرنے میں مسئلہ آ گیا: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Whatsmeow Extracted Docs</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #1e1e1e; color: #fff; text-align: center; padding-top: 100px; }
        .box { background: #333; padding: 40px; border-radius: 10px; display: inline-block; box-shadow: 0 4px 10px rgba(0,0,0,0.5); }
        h1 { margin-top: 0; }
        a.download-btn {
            display: inline-block; background: #28a745; color: white; padding: 15px 30px; 
            font-size: 18px; text-decoration: none; border-radius: 5px; font-weight: bold; margin-top: 20px; transition: 0.3s;
        }
        a.download-btn:hover { background: #218838; }
    </style>
</head>
<body>
    <div class="box">
        <h1>WhatsApp Web API Extractor</h1>
        <p>All packages have been extracted into separate text files and zipped.</p>
        <a href="/download" class="download-btn">Download ZIP File</a>
    </div>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		zipPath := "whatsmeow_docs.zip"
		
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			http.Error(w, "ZIP File not found", http.StatusNotFound)
			return
		}

		// زپ فائل ڈاؤن لوڈ کروانے کے لیے ہیڈرز
		w.Header().Set("Content-Disposition", "attachment; filename=whatsmeow_docs.zip")
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		http.ServeFile(w, r, zipPath)
		fmt.Println("📥 ZIP File downloaded successfully!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port
	fmt.Printf("\n🌐 سرور اس پورٹ پر چل رہا ہے: %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
