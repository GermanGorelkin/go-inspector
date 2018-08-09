package main

import (
	"github.com/germangorelkin/go-inspector/inspector"
	"log"
	"net/url"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set.")
	}
	instance := os.Getenv("INSTANCE")
	if instance == "" {
		log.Fatal("INSTANCE is not set.")
	}

	log.Print(apiKey)
	log.Print(instance)
	//return

	inst, err := url.Parse(instance)
	if err != nil {
		log.Fatal(err)
	}
	cfg := inspector.ClintConf{
		APIKey:   apiKey,
		Instance: inst,
	}
	c := inspector.NewClient(cfg)

	//f, err := os.Open("C://Users//gg//Desktop//hotfield_4e9a6589b0c0ad6a1533543806228.jpg")
	//if err != nil {
	//	log.Panic(err)
	//}
	//defer f.Close()

	//Uploads
	//u, err := c.Image.Uploads(f, f.Name())
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("%+v", u)

	// Recognize
	//recreq := &inspector.RecognizeRequest{
	//	Images:      []int{4973284},
	//	ReportTypes: []string{inspector.ReportTypeFACING_COUNT, inspector.ReportTypePRICE_TAGS},
	//}
	//recres, err := c.Recognize.Recognize(recreq)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("%+v", recres)

	// Report
	//report, err := c.Report.GetReport(14621)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//r, err := c.Report.ToPriceTags(report.Json)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("%+v\n", r)

	//report, err := c.Report.GetReport(14622)
	//if err != nil {
	//	log.Println(err)
	//}
	//r, err := c.Report.ToFacingCount(report.Json)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("%+v\n", r)

	pag, err := c.Sku.GetSKU(0, 10)
	if err != nil {
		log.Println(err)
	}
	//log.Printf("%+v\n", pag)
	sku, err := c.Sku.ToSku(pag.Results)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", sku)
}
