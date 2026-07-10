package photo

import "time"

type Photo struct {
	ID               int64
	Filename         string
	Note             string
	TenantID         int64
	InspectionActID  *int64
	ReplacementActID *int64
	OriginalFilename string
	ContentType      string
	SizeBytes        int64
	UploadedBy       *int64
	CreatedAt        time.Time
}

// Summary is the projection exposed to the act domain (via the PhotoLister
// port) to list attachments and embed them in the generated PDF.
type Summary struct {
	ID               int64
	OriginalFilename string
	Note             string
	FilePath         string
	CreatedAt        time.Time
}
