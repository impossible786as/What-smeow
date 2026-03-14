package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	// یہ بلینک امپورٹ گو کو بتائے گا کہ whatsmeow کو ڈیلیٹ نہیں کرنا!
	_ "go.mau.fi/whatsmeow"
)

func generateDocs() error {
	fileName := "whatsmeow_full_functions.txt"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("فائل بنانے میں مسئلہ: %v", err)
	}
	defer file.Close()

	fmt.Println("\n🚀 ----------------------------------------------------")
	fmt.Println("🚀 STEP 1: واٹس میو کے تمام پیکجز تلاش کیے جا رہے ہیں...")
	fmt.Println("🚀 ----------------------------------------------------")

	cmdList := exec.Command("go", "list", "go.mau.fi/whatsmeow/...")
	var listErr bytes.Buffer
	cmdList.Stderr = &listErr
	outputList, err := cmdList.Output()

	if err != nil {
		fmt.Printf("❌ ایرر 'go list' کمانڈ چلانے میں: %v\n", err)
		fmt.Printf("❌ اندرونی ایرر (STDERR): %s\n", listErr.String())
		return fmt.Errorf("پیکجز فائنڈ نہیں ہوئے")
	}

	packagesStr := strings.TrimSpace(string(outputList))
	if packagesStr == "" {
		fmt.Println("⚠️ وارننگ: کوئی پیکج نہیں ملا، لسٹ بالکل خالی ہے!")
		return fmt.Errorf("empty package list")
	}

	packages := strings.Split(packagesStr, "\n")
	fmt.Printf("✅ زبردست! ٹوٹل %d پیکجز مل گئے ہیں۔\n\n", len(packages))

	fmt.Println("🚀 ----------------------------------------------------")
	fmt.Println("🚀 STEP 2: ہر پیکج کے فنکشنز ایکسٹریکٹ کیے جا رہے ہیں...")
	fmt.Println("🚀 ----------------------------------------------------\n")

	for i, pkg := range packages {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			continue
		}

		fmt.Printf("⏳ [%d/%d] پروسیسنگ ہو رہی ہے: %s\n", i+1, len(packages), pkg)

		cmdDoc := exec.Command("go", "doc", "-all", pkg)
		var docErr bytes.Buffer
		cmdDoc.Stderr = &docErr
		outputDoc, err := cmdDoc.Output()

		if err != nil {
			fmt.Printf("  ❌ فیل ہو گیا! ایرر: %v\n", err)
			continue
		}

		if len(outputDoc) == 0 {
			fmt.Printf("  ⚠️ کوئی پبلک فنکشن یا سٹرکچر نہیں ملا۔\n")
		} else {
			fmt.Printf("  ✅ کامیابی! اس پیکج سے %d بائٹس کا ڈیٹا نکالا گیا۔\n", len(outputDoc))
			
			file.WriteString("========================================================\n")
			file.WriteString(fmt.Sprintf("PACKAGE: %s\n", pkg))
			file.WriteString("========================================================\n\n")
			file.Write(outputDoc)
			file.WriteString("\n\n\n")
		}
	}

	fmt.Println("\n🎉 سب کچھ مکمل ہو گیا! فائل ریڈی ہے۔")
	return nil
}

func main() {
	fmt.Println("=================================================")
	fmt.Println("   WHATSMEOW EXTRACTOR SERVER STARTING...        ")
	fmt.Println("=================================================")

	if err := generateDocs(); err != nil {
		fmt.Printf("\n🚨 خطرہ: ڈاکیومنٹیشن جنریٹ کرنے میں مسئلہ آ گیا: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("whatsmeow_full_functions.txt")
		
		if err != nil || len(content) == 0 {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `
				<body style="background:#1e1e1e; color:red; font-family:Arial; text-align:center; padding-top:50px;">
					<h2>🚨 Error: File is empty or not generated yet!</h2>
					<p>پلیز Railway کے لاگز چیک کریں کہ کیا مسئلہ آیا ہے۔</p>
				</body>
			`)
			return
		}

		html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Whatsmeow Extracted Docs</title>
    <style>
        body { font-family: Consolas, monospace; background-color: #1e1e1e; color: #d4d4d4; margin: 0; padding: 15px; }
        .header { 
            display: flex; justify-content: space-between; align-items: center; 
            background: #333; padding: 15px; border-radius: 8px; 
            margin-bottom: 20px; position: sticky; top: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
        h2 { margin: 0; color: #fff; font-size: 16px; font-family: Arial, sans-serif; }
        button { 
            background: #007bff; color: white; border: none; 
            padding: 10px 20px; font-size: 14px; border-radius: 5px; 
            cursor: pointer; font-weight: bold; transition: 0.3s;
        }
        button:hover { background: #0056b3; }
        pre { 
            background: #2d2d2d; padding: 15px; border-radius: 8px; 
            overflow-x: auto; white-space: pre-wrap; word-wrap: break-word; 
            font-size: 13px; line-height: 1.5;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>Whatsmeow Functions</h2>
        <button id="copyBtn" onclick="copyDoc()">Copy All Text</button>
    </div>
    <pre id="docContent">%s</pre>

    <script>
        function copyDoc() {
            const text = document.getElementById("docContent").innerText;
            navigator.clipboard.writeText(text).then(() => {
                const btn = document.getElementById('copyBtn');
                btn.innerText = "Copied! ✔";
                btn.style.background = "#28a745";
                setTimeout(() => {
                    btn.innerText = "Copy All Text";
                    btn.style.background = "#007bff";
                }, 3000);
            }).catch(err => {
                alert("Copy failed! Please select manually.");
            });
        }
    </script>
</body>
</html>`

		safeContent := template.HTMLEscapeString(string(content))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, html, safeContent)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port
	fmt.Printf("\n🌐 سرور اس پورٹ پر چل رہا ہے: %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
