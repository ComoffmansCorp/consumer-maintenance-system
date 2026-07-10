package act

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/signintech/gopdf"
)

//go:embed assets/fonts/DejaVuSans.ttf
var pdfFontRegularData []byte

//go:embed assets/fonts/DejaVuSans-Bold.ttf
var pdfFontBoldData []byte

const (
	pdfFontRegular = "DejaVuSans"
	pdfFontBold    = "DejaVuSans-Bold"

	pdfMargin       = 18.0
	pdfPageWidth    = 210.0
	pdfPageHeight   = 297.0
	pdfContentWidth = pdfPageWidth - 2*pdfMargin
	pdfBottomLimit  = pdfPageHeight - pdfMargin
)

// pdfDocument wraps gopdf with the fonts this domain always needs and a
// simple flowing-Y layout helper, so the act generation logic below reads
// top to bottom like the document it produces.
type pdfDocument struct {
	pdf *gopdf.GoPdf
	y   float64
}

func newPDFDocument() (*pdfDocument, error) {
	pdf := &gopdf.GoPdf{}
	// Unit: UnitMM makes every coordinate/size below (margins, cell rects,
	// line endpoints) interpreted as millimeters; PageSizeA4 carries its own
	// point-based override so the physical page size stays correct.
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, Unit: gopdf.UnitMM})
	if err := pdf.AddTTFFontData(pdfFontRegular, pdfFontRegularData); err != nil {
		return nil, fmt.Errorf("load regular font: %w", err)
	}
	if err := pdf.AddTTFFontData(pdfFontBold, pdfFontBoldData); err != nil {
		return nil, fmt.Errorf("load bold font: %w", err)
	}
	pdf.AddPage()
	doc := &pdfDocument{pdf: pdf, y: pdfMargin}
	return doc, nil
}

func (d *pdfDocument) ensureSpace(height float64) {
	if d.y+height > pdfBottomLimit {
		d.pdf.AddPage()
		d.y = pdfMargin
	}
}

func (d *pdfDocument) gap(h float64) {
	d.y += h
}

func (d *pdfDocument) text(style string, size float64, text string) {
	_ = d.pdf.SetFont(style, "", size)
	d.ensureSpace(size / 2)
	d.pdf.SetXY(pdfMargin, d.y)
	_ = d.pdf.Cell(nil, text)
	d.y += size/2 + 2
}

func (d *pdfDocument) heading(text string) {
	d.text(pdfFontBold, 16, text)
}

func (d *pdfDocument) subheading(text string) {
	d.gap(2)
	d.text(pdfFontBold, 12, text)
}

// row prints a "label: value" line with the label bold and the value in the
// monospace-flavored regular weight (the design system reserves true
// monospace for the web UI; here regular weight keeps the PDF legible).
func (d *pdfDocument) row(label, value string) {
	if value == "" {
		value = "—"
	}
	_ = d.pdf.SetFont(pdfFontBold, "", 10)
	d.ensureSpace(6)
	d.pdf.SetXY(pdfMargin, d.y)
	_ = d.pdf.Cell(&gopdf.Rect{W: 55, H: 6}, label+":")

	_ = d.pdf.SetFont(pdfFontRegular, "", 10)
	d.pdf.SetXY(pdfMargin+55, d.y)
	_ = d.pdf.MultiCell(&gopdf.Rect{W: pdfContentWidth - 55, H: 6}, value)
	d.y += 6
}

func (d *pdfDocument) hline() {
	d.gap(2)
	d.ensureSpace(2)
	d.pdf.SetLineWidth(0.3)
	d.pdf.SetStrokeColor(200, 200, 200)
	d.pdf.Line(pdfMargin, d.y, pdfMargin+pdfContentWidth, d.y)
	d.gap(4)
}

func (d *pdfDocument) bytes() []byte {
	return d.pdf.GetBytesPdf()
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("02.01.2006")
}

func formatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

func formatReadings(v *float64) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%.3f", *v)
}

func formatIntPtr(v *int32) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

func translateInspectionType(t InspectionType) string {
	switch t {
	case InspectionScheduled:
		return "Плановый"
	case InspectionUnscheduled:
		return "Внеплановый"
	default:
		return string(t)
	}
}

func translateMeterType(t string) string {
	switch t {
	case "SINGLE_PHASE":
		return "Однофазный"
	case "THREE_PHASE_DIRECT":
		return "Трёхфазный прямого включения"
	case "THREE_PHASE_TRANSFORMER":
		return "Трёхфазный трансформаторного включения"
	default:
		return t
	}
}

func translateSealState(s string) string {
	switch s {
	case "INTACT":
		return "Не нарушена"
	case "BROKEN":
		return "Нарушена"
	case "MISSING":
		return "Отсутствует"
	default:
		return "—"
	}
}

func (d *pdfDocument) letterhead(tenant TenantInfo, actKind, actNumber string, generatedAt time.Time) {
	if tenant.Name != "" {
		d.heading(tenant.Name)
	} else {
		d.heading("Организация обслуживания приборов учёта")
	}
	d.text(pdfFontRegular, 9, "Акт № "+actNumber+" сформирован "+formatDateTime(generatedAt))
	d.hline()
	d.subheading(actKind)
}

func (d *pdfDocument) metersTable(meters []MeterInfo) {
	if len(meters) == 0 {
		d.text(pdfFontRegular, 10, "Приборы учёта не зафиксированы.")
		return
	}
	headers := []string{"Тип", "Серийный номер", "Год вып.", "Поверка", "Пломба", "Коэфф."}
	widths := []float64{40, 40, 20, 25, 27, 22}

	d.ensureSpace(7)
	_ = d.pdf.SetFont(pdfFontBold, "", 9)
	x := pdfMargin
	for i, h := range headers {
		d.pdf.SetXY(x, d.y)
		_ = d.pdf.CellWithOption(&gopdf.Rect{W: widths[i], H: 7}, h, gopdf.CellOption{Border: gopdf.AllBorders})
		x += widths[i]
	}
	d.y += 7

	_ = d.pdf.SetFont(pdfFontRegular, "", 9)
	for _, m := range meters {
		d.ensureSpace(7)
		x = pdfMargin
		values := []string{
			translateMeterType(m.Type),
			m.SerialNumber,
			formatIntPtr(m.ManufactureYear),
			formatDate(m.VerificationDate),
			translateSealState(m.SealState),
			formatIntPtr(m.TransformationRatio),
		}
		for i, v := range values {
			d.pdf.SetXY(x, d.y)
			_ = d.pdf.CellWithOption(&gopdf.Rect{W: widths[i], H: 7}, v, gopdf.CellOption{Border: gopdf.AllBorders})
			x += widths[i]
		}
		d.y += 7
	}
}

func (d *pdfDocument) photosAppendix(photos []PhotoInfo) {
	if len(photos) == 0 {
		return
	}
	d.subheading(fmt.Sprintf("Фотофиксация (%d)", len(photos)))
	for _, p := range photos {
		d.embedPhoto(p)
	}
}

func (d *pdfDocument) embedPhoto(p PhotoInfo) {
	const maxW, maxH = 80.0, 60.0

	caption := p.OriginalFilename
	if p.Note != "" {
		caption += " — " + p.Note
	}

	if p.FilePath == "" {
		d.text(pdfFontRegular, 9, "• "+caption+" (файл недоступен)")
		return
	}
	if _, err := os.Stat(p.FilePath); err != nil {
		d.text(pdfFontRegular, 9, "• "+caption+" (файл недоступен)")
		return
	}

	d.ensureSpace(maxH + 8)
	if err := d.pdf.Image(p.FilePath, pdfMargin, d.y, &gopdf.Rect{W: maxW, H: maxH}); err != nil {
		d.text(pdfFontRegular, 9, "• "+caption+" (не удалось встроить изображение)")
		return
	}
	d.y += maxH + 2
	d.text(pdfFontRegular, 8, caption)
	d.gap(2)
}
