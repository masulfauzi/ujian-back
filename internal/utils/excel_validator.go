package utils

import (
	"strings"
)

// ValidateSoalRow memvalidasi satu row dari excel
func ValidateSoalRow(row *ExcelSoalRow) []string {
	var errors []string

	// Validate no_soal
	if row.NoSoal <= 0 {
		errors = append(errors, "no_soal harus lebih dari 0")
	}

	// Validate soal
	if strings.TrimSpace(row.Soal) == "" {
		errors = append(errors, "soal tidak boleh kosong")
	}
	if len(row.Soal) > 5000 {
		errors = append(errors, "soal max 5000 karakter")
	}

	// Validate opsi_a
	if strings.TrimSpace(row.OpsiA) == "" {
		errors = append(errors, "opsi_a tidak boleh kosong")
	}
	if len(row.OpsiA) > 5000 {
		errors = append(errors, "opsi_a max 5000 karakter")
	}

	// Validate opsi_b
	if strings.TrimSpace(row.OpsiB) == "" {
		errors = append(errors, "opsi_b tidak boleh kosong")
	}
	if len(row.OpsiB) > 5000 {
		errors = append(errors, "opsi_b max 5000 karakter")
	}

	// Validate opsi_c
	if strings.TrimSpace(row.OpsiC) == "" {
		errors = append(errors, "opsi_c tidak boleh kosong")
	}
	if len(row.OpsiC) > 5000 {
		errors = append(errors, "opsi_c max 5000 karakter")
	}

	// Validate opsi_d
	if strings.TrimSpace(row.OpsiD) == "" {
		errors = append(errors, "opsi_d tidak boleh kosong")
	}
	if len(row.OpsiD) > 5000 {
		errors = append(errors, "opsi_d max 5000 karakter")
	}

	// Validate opsi_e
	if strings.TrimSpace(row.OpsiE) == "" {
		errors = append(errors, "opsi_e tidak boleh kosong")
	}
	if len(row.OpsiE) > 5000 {
		errors = append(errors, "opsi_e max 5000 karakter")
	}

	// Validate kunci
	kunci := strings.TrimSpace(row.Kunci)
	if kunci == "" {
		errors = append(errors, "kunci tidak boleh kosong")
	}
	if !IsValidKunci(kunci) {
		errors = append(errors, "kunci harus berupa A, B, C, D, atau E")
	}

	return errors
}

// IsValidKunci checks if kunci is valid (A-E)
func IsValidKunci(k string) bool {
	k = strings.ToUpper(strings.TrimSpace(k))
	return k == "A" || k == "B" || k == "C" || k == "D" || k == "E"
}
