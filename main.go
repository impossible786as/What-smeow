package main

import (
	"fmt"
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
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div style="font-family: Arial, sans-serif; text-align: center; margin-top: 50px;">
				<h1>WhatsApp Web API Extractor</h1>
				<p>All packages have been extracted successfully.</p>
				<br><br>
				<a href="/download" style="font-size: 18px; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 8px;">
					Download whatsmeow_full_functions.txt
				</a>
			</div>
		`)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Download requested by browser...")
		
		// پوری فائل کو میموری میں ریڈ کرنا
		fileBytes, err := os.ReadFile("whatsmeow_full_functions.txt")
		if err != nil {
			http.Error(w, "File not found or still generating", http.StatusInternalServerError)
			log.Printf("Download error: %v", err)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=whatsmeow_full_functions.txt")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// براؤزر اور ریلوے سرور کو ایگزیکٹ سائز بتانا تاکہ کنکشن ڈراپ نہ ہو
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))

		// ڈائریکٹ بائٹس سینڈ کرنا
		w.Write(fileBytes)
		fmt.Println("Download completed successfully!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port
	fmt.Printf("Server is running on %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
