package utils

import (
	"strings"
)

// ValidateSoalRow memvalidasi satu row dari excel. Jenis & kunci yang valid
// ditulis kembali ke row dalam bentuk ternormalisasi (mis. kunci pgk "C,A" → "AC").
func ValidateSoalRow(row *ExcelSoalRow) []string {
	var errors []string

	jenis, err := NormalizeJenis(row.Jenis)
	if err != nil {
		return append(errors, err.Error())
	}
	row.Jenis = jenis

	if strings.TrimSpace(row.Kunci) == "" {
		return append(errors, "kunci tidak boleh kosong")
	}
	kunci, err := NormalizeKunci(jenis, row.Kunci)
	if err != nil {
		return append(errors, err.Error())
	}
	row.Kunci = kunci

	opsi := []string{row.OpsiA, row.OpsiB, row.OpsiC, row.OpsiD, row.OpsiE}
	adaGambar := []bool{row.GambarA != "", row.GambarB != "", row.GambarC != "", row.GambarD != "", row.GambarE != ""}
	if err := ValidateOpsiSoal(jenis, kunci, opsi, adaGambar); err != nil {
		errors = append(errors, err.Error())
	}

	return errors
}

// IsValidKunci checks if kunci is valid (A-E)
func IsValidKunci(k string) bool {
	k = strings.ToUpper(strings.TrimSpace(k))
	return k == "A" || k == "B" || k == "C" || k == "D" || k == "E"
}
