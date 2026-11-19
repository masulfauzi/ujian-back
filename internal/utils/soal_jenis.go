package utils

import (
	"errors"
	"sort"
	"strings"

	"backend/internal/constants"
)

// NormalizeJenis mengembalikan jenis soal yang valid (pg, pgk, bs). Kosong dianggap pg.
func NormalizeJenis(j string) (string, error) {
	j = strings.ToLower(strings.TrimSpace(j))
	switch j {
	case "":
		return constants.JenisSoalPG, nil
	case constants.JenisSoalPG, constants.JenisSoalPGK, constants.JenisSoalBS:
		return j, nil
	default:
		return "", errors.New("jenis soal harus salah satu dari: pg, pgk, bs")
	}
}

// cleanHuruf uppercase dan membuang spasi/koma, mis. "a, c" → "AC".
func cleanHuruf(s string) string {
	s = strings.ToUpper(s)
	return strings.NewReplacer(" ", "", ",", "", ";", "").Replace(s)
}

// sortUnikOpsi memvalidasi huruf A-E, membuang duplikat, lalu mengurutkan.
func sortUnikOpsi(s string) (string, bool) {
	seen := map[rune]bool{}
	var out []string
	for _, r := range s {
		if r < 'A' || r > 'E' {
			return "", false
		}
		if !seen[r] {
			seen[r] = true
			out = append(out, string(r))
		}
	}
	sort.Strings(out)
	return strings.Join(out, ""), true
}

func normalizeBS(s string) (string, error) {
	if len(s) != 3 || strings.Trim(s, "BS") != "" {
		return "", errors.New("soal benar/salah harus berisi tepat 3 karakter B atau S, mis. BSB")
	}
	return s, nil
}

// NormalizeKunci memvalidasi dan menormalkan kunci sesuai jenis soal.
func NormalizeKunci(jenis, kunci string) (string, error) {
	k := cleanHuruf(kunci)
	switch jenis {
	case constants.JenisSoalPGK:
		norm, ok := sortUnikOpsi(k)
		if !ok || len(norm) < 2 {
			return "", errors.New("kunci pilihan ganda kompleks harus minimal 2 huruf berbeda dari A-E, mis. AC")
		}
		return norm, nil
	case constants.JenisSoalBS:
		return normalizeBS(k)
	default:
		if !IsValidKunci(k) || len(k) != 1 {
			return "", errors.New("kunci harus A, B, C, D, atau E")
		}
		return k, nil
	}
}

// NormalizeJawaban memvalidasi dan menormalkan jawaban peserta sesuai jenis soal.
func NormalizeJawaban(jenis, jawaban string) (string, error) {
	j := cleanHuruf(jawaban)
	switch jenis {
	case constants.JenisSoalPGK:
		norm, ok := sortUnikOpsi(j)
		if !ok || norm == "" {
			return "", errors.New("jawaban pilihan ganda kompleks harus berisi huruf A-E, mis. AC")
		}
		return norm, nil
	case constants.JenisSoalBS:
		return normalizeBS(j)
	default:
		if !IsValidKunci(j) || len(j) != 1 {
			return "", errors.New("jawaban harus salah satu dari: A, B, C, D, E")
		}
		return j, nil
	}
}

// HitungSkor menghitung poin satu jawaban yang sudah dinormalisasi.
// pg/pgk: sama persis dengan kunci → 1, selain itu 0.
// bs: 0-1 pernyataan benar → 0, 2 benar → 1, 3 benar → 2. isBenar = 1 hanya jika poin penuh.
func HitungSkor(jenis, jawaban, kunci string) (skor, isBenar int) {
	if jenis == constants.JenisSoalBS {
		cocok := 0
		for i := 0; i < len(jawaban) && i < len(kunci); i++ {
			if jawaban[i] == kunci[i] {
				cocok++
			}
		}
		if cocok > 1 {
			skor = cocok - 1
		}
		if cocok == 3 {
			isBenar = 1
		}
		return skor, isBenar
	}
	if jawaban == kunci {
		return 1, 1
	}
	return 0, 0
}

// MapHurufOpsi memetakan tiap huruf opsi antara posisi asli dan posisi acak.
// order[i] = huruf asli yang ditampilkan di posisi i.
// reverse=false: asli → tampilan (untuk kunci yang dikirim ke frontend).
// reverse=true: tampilan → asli (untuk mengevaluasi jawaban peserta).
func MapHurufOpsi(huruf string, order []string, reverse bool) string {
	var out []string
	for _, r := range huruf {
		h := string(r)
		if reverse {
			idx := int(r - 'A')
			if idx >= 0 && idx < len(order) {
				h = order[idx]
			}
		} else {
			for i, key := range order {
				if key == h {
					h = string(rune('A' + i))
					break
				}
			}
		}
		out = append(out, h)
	}
	sort.Strings(out)
	return strings.Join(out, "")
}

// ValidateOpsiSoal memastikan opsi sesuai jenis soal. opsi & adaGambar berurutan A-E.
// bs: pernyataan 1-3 (opsi A-C) wajib terisi. pgk: setiap huruf kunci wajib punya opsi/gambar.
func ValidateOpsiSoal(jenis, kunci string, opsi []string, adaGambar []bool) error {
	terisi := func(i int) bool {
		return strings.TrimSpace(opsi[i]) != "" || adaGambar[i]
	}
	switch jenis {
	case constants.JenisSoalBS:
		for i := 0; i < 3; i++ {
			if !terisi(i) {
				return errors.New("soal benar/salah wajib berisi 3 pernyataan di opsi_a, opsi_b, dan opsi_c")
			}
		}
	case constants.JenisSoalPGK:
		for _, r := range kunci {
			if !terisi(int(r - 'A')) {
				return errors.New("opsi " + string(r) + " pada kunci pilihan ganda kompleks tidak boleh kosong")
			}
		}
	}
	return nil
}
