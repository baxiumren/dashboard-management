package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"
)

func buildFuncMap() template.FuncMap {
	return template.FuncMap{
		"inc": func(i int) int { return i + 1 },
		"title": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"titleTipe": func(s string) string {
			m := map[string]string{
				"buat":       "Fanpage Dibuat",
				"beli":       "Beli Akun",
				"checkpoint": "Checkpoint",
				"suspend":    "Suspend",
				"kunci":      "Kunci",
				"selfie":     "Selfie / Verifikasi",
				"disabled":   "Disabled",
				"nonaktif":   "Nonaktif",
				"banned":     "Banned",
				"pulih":      "Pulih / Aktif Kembali",
				"lainnya":    "Lainnya",
			}
			if v, ok := m[s]; ok {
				return v
			}
			return s
		},
		"tipeColor": func(s string) string {
			m := map[string]string{
				"buat":       "background:#3b82f6;border-color:#3b82f6",
				"beli":       "background:#16a34a;border-color:#16a34a",
				"checkpoint": "background:#b45309;border-color:#b45309",
				"suspend":    "background:#dc2626;border-color:#dc2626",
				"kunci":      "background:#7c3aed;border-color:#7c3aed",
				"selfie":     "background:#0891b2;border-color:#0891b2",
				"disabled":   "background:#b45309;border-color:#b45309",
				"nonaktif":   "background:#b91c1c;border-color:#b91c1c",
				"banned":     "background:#6b7280;border-color:#6b7280",
				"pulih":      "background:#16a34a;border-color:#16a34a",
				"lainnya":    "background:#64748b;border-color:#64748b",
			}
			if v, ok := m[s]; ok {
				return v
			}
			return "background:#64748b;border-color:#64748b"
		},
		"jsonStr": func(s string) template.JS {
			b, _ := json.Marshal(s)
			return template.JS(b)
		},
		"formatRp": func(n int) string {
			return fmt.Sprintf("Rp %s", formatNumber(n))
		},
	}
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

// cloneParse clones the base template then parses additional page-specific files.
// Each page gets its own isolated template set so "content" blocks don't conflict.
func cloneParse(base *template.Template, files ...string) *template.Template {
	t, err := base.Clone()
	if err != nil {
		log.Fatalf("Template clone error: %v", err)
	}
	t, err = t.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Template parse error: %v", err)
	}
	return t
}
