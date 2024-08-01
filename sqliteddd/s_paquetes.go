package sqliteddd

import "fmt"

func (s *Repositorio) ExistePaquete(paqueteID int, nombre string) bool {
	existe := false
	err := s.db.QueryRow("SELECT 1 FROM paquetes WHERE paquete_id = ? OR nombre = ?", paqueteID, nombre).Scan(&existe)
	if err != nil {
		fmt.Println(err)
	}
	return existe
}
