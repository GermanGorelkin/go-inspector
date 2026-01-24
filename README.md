# go inspector
Go API for Inspector Cloud (tested on Go 1.24)

## Project Status
Beta 0.1

## Install
``
go get github.com/germangorelkin/go-inspector/inspector
``

## Quickstart

### Init client
```go
apiKey := ""
instance := ""

// logrus.SetLevel(logrus.DebugLevel)

inst, err := url.Parse(instance)
if err != nil {
	log.Fatal(err)
}
cfg := inspector.ClintConf{
	APIKey:   apiKey,
	Instance: inst,
}
c := inspector.NewClient(cfg)
```

### Add Visit
```go
v := &inspector.Visit{}
v, err = c.Visit.AddVisit(v)
if err != nil {
	log.Println(err)
}
log.Printf("%+v\n", v)
```

### Upload Image
```go
f, err := os.Open("")
if err != nil {
	log.Panic(err)
}
defer f.Close()

u, err := c.Image.Upload(f, f.Name())
if err != nil {
	log.Println(err)
}
log.Printf("%+v", u)
```
or
```go
u, err := c.Image.UploadByURL("")
if err != nil {
	logrus.Errorf("%+v\n", err)
}
log.Printf("result: %+v", u)
```

### Recognize
```go
imageID = 111111
recreq := &inspector.RecognizeRequest{
	Images:      []int{imageID},
	ReportTypes: []string{inspector.ReportTypeFACING_COUNT, inspector.ReportTypePRICE_TAGS},
}
recres, err := c.Recognize.Recognize(recreq)
if err != nil {
	log.Println(err)
}
log.Printf("%+v", recres)
```

### Get Report
```go
// PriceTags report
reportID = 1111
report, err := c.Report.GetReport(reportID)
if err != nil {
	log.Println(err)
}
	
r, err := c.Report.ToPriceTags(report.Json)
if err != nil {
	log.Println(err)
}
log.Printf("%+v\n", r)

// FacingCount report
reportID = 1112
report, err := c.Report.GetReport(reportID)
if err != nil {
	log.Println(err)
}
r, err := c.Report.ToFacingCount(report.Json)
if err != nil {
	log.Println(err)
}
log.Printf("%+v\n", r)
```

### Fetch SKU
```go
pag, err := c.Sku.GetSKU(0, 10)
if err != nil {
	log.Println(err)
}
	
sku, err := c.Sku.ToSku(pag.Results)
if err != nil {
	log.Println(err)
}
log.Printf("%+v\n", sku)
```
