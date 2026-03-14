package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func generateDocs() error {
	fileName := "whatsmeow_full_functions.txt"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	cmdList := exec.Command("go", "list", "go.mau.fi/whatsmeow/...")
	outputList, err := cmdList.Output()
	if err != nil {
		return fmt.Errorf("failed to list packages: %v", err)
	}

	packages := strings.Split(strings.TrimSpace(string(outputList)), "\n")

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}
		
		cmdDoc := exec.Command("go", "doc", "-all", pkg)
		outputDoc, err := cmdDoc.CombinedOutput()
		if err != nil {
			continue
		}
		
		file.WriteString("========================================================\n")
		file.WriteString(fmt.Sprintf("PACKAGE: %s\n", pkg))
		file.WriteString("========================================================\n\n")
		file.Write(outputDoc)
		file.WriteString("\n\n\n")
	}

	return nil
}

func main() {
	fmt.Println("Starting Railway Server...")
	fmt.Println("Hunting down EVERY SINGLE package inside whatsmeow...")
	
	if err := generateDocs(); err != nil {
		log.Fatalf("Error generating docs: %v", err)
	}
	fmt.Println("Full comprehensive documentation generated successfully!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// فائل کو ریڈ کرنا
		content, err := os.ReadFile("whatsmeow_full_functions.txt")
		if err != nil {
			http.Error(w, "File not found or still generating", http.StatusInternalServerError)
			return
		}

		// HTML اور JavaScript کا کوڈ تاکہ سکرین پر نظر آئے اور کاپی ہو سکے
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
                
                // 3 سیکنڈ بعد بٹن دوبارہ نارمل ہو جائے گا
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

		// ٹیکسٹ کو HTML کے لیے Safe بنانا تاکہ < > والے ٹیگز خراب نہ ہوں
		safeContent := template.HTMLEscapeString(string(content))
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, html, safeContent)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port
	fmt.Printf("Server is running on %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
