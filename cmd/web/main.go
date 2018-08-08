package main

import (
	"fmt"
	"github.com/germangorelkin/go-inspector/inspector"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	apiKey   = ""
	instance = ""
)

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html"))

	if r.Method == http.MethodGet {
		tmpl.Execute(w, "")
	} else if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		inst, err := url.Parse(instance)
		if err != nil {
			log.Panic(err)
		}
		cfg := inspector.ClintConf{
			APIKey:   apiKey,
			Instance: inst,
		}
		c := inspector.NewClient(cfg)

		//Uploads
		u, err := c.Image.Uploads(file, handler.Filename)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%+v\n", u)

		// Recognize
		recreq := &inspector.RecognizeRequest{
			Images:      []int{u.ID},
			ReportTypes: []string{inspector.ReportTypeFACING_COUNT, inspector.ReportTypePRICE_TAGS},
		}
		recres, err := c.Recognize.Recognize(recreq)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v", recres)

		// Report
		idPriceTags := recres.Reports[inspector.ReportTypePRICE_TAGS]
		var report *inspector.Report
		for i := 0; i < 3; i++ {
			time.Sleep(2 * time.Second)
			report, err = c.Report.GetReport(idPriceTags)
			if err != nil {
				log.Println(err)
			}
			if report.Status == inspector.ReportStatusREADY {
				continue
			}
			time.Sleep(2 * time.Second)
		}
		priceTags, err := c.Report.ToPriceTags(report.Json)
		if err != nil {
			log.Println(err)
		}
		//log.Printf("%+v\n\n", priceTags)

		idFacingCount := recres.Reports[inspector.ReportTypeFACING_COUNT]
		for i := 0; i < 3; i++ {
			//time.Sleep(1 * time.Second)
			report, err = c.Report.GetReport(idFacingCount)
			if err != nil {
				log.Println(err)
			}
			if report.Status == inspector.ReportStatusREADY {
				continue
			}
			time.Sleep(2 * time.Second)
		}
		facingCount, err := c.Report.ToFacingCount(report.Json)
		if err != nil {
			log.Println(err)
		}
		//log.Printf("%+v\n\n", facingCount)

		var v = struct {
			Visit                 int
			ReportFacingCountJson []inspector.ReportFacingCountJson
			ReportPriceTagsJson   []inspector.ReportPriceTagsJson
		}{
			Visit: report.Visit,
			ReportFacingCountJson: facingCount,
			ReportPriceTagsJson:   priceTags,
		}
		tmpl.Execute(w, v)
	}
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set.")
	}
	instance := os.Getenv("INSTANCE")
	if instance == "" {
		log.Fatal("INSTANCE is not set.")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not set.")
	}
	addr := ":" + port

	http.HandleFunc("/", handler)
	log.Printf("The service(%s) is ready to listen and serve.", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
