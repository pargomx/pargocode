package plantillas

import (
	"fmt"
	"time"
)

// Convirte duración en hh:mm.
// Ejemplo: 1h1m45s -> "01:02"
func fmtHorasMinutos(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
