package inspector

import "time"

// HTTP methods
const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"
	methodPATCH  = "PATCH"
)

// Endpoint paths for IC API
const (
	// Image endpoints
	endpointUploads      = "uploads/"
	endpointUploadsByURL = "uploads/upload_by_url/"

	// Recognition endpoints
	endpointRecognize        = "recognize/"
	endpointRecognitionError = "recognition_error/"

	// Report endpoints
	endpointReports = "reports/%d/" // formatted with report ID

	// SKU endpoints
	endpointSKU = "sku/"

	// Visit endpoints
	endpointVisits = "visits/"
)

// Default timeouts and intervals
const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 30 * time.Second

	// DefaultPollingInterval is the default interval for polling report status
	DefaultPollingInterval = 2 * time.Second

	// DefaultPollingTimeout is the default overall timeout for polling
	DefaultPollingTimeout = 60 * time.Second
)

// Authentication scheme
const (
	authSchemeToken = "Token"
)

// Default values
const (
	// DefaultPageSize is the default page size for paginated requests
	DefaultPageSize = 100

	// MaxPaginationPages is the maximum number of pages to fetch
	// to prevent infinite loops in pagination
	MaxPaginationPages = 1000
)

// HTTP header names
const (
	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"
)

// Multipart form field names
const (
	formFieldFile = "file"
)

// Content types
const (
	contentTypeMultipartFormData = "multipart/form-data"
	contentTypeJSON              = "application/json"
)
