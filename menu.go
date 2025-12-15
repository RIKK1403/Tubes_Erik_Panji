package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Code  int
	Name  string
	Price int64 // rupiah
}

type Line struct {
	Item      Item
	Qty       int
	LineTotal int64
}

var menu = []Item{
	{1, "Nasi Goreng", 22000},
	{2, "Mie Goreng", 18000},
	{3, "Ayam Goreng", 25001},
	{4, "Es Teh", 6000},
	{5, "Air Mineral", 5000},
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	cart := make(map[int]int)

	fmt.Println("=== Selamat Datang di ERPE ===")

	// ================= INPUT BARANG =================
	for {
		printMenu()
		fmt.Print("Masukkan kode item (0 untuk Pembayaran): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		kode, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Input harus berupa angka.")
			continue
		}

		if kode == 0 {
			break
		}

		item, ok := findItemByCode(kode)
		if !ok {
			fmt.Println("Kode item tidak ditemukan.")
			continue
		}

		fmt.Printf("Masukkan jumlah %s: ", item.Name)
		qtyInput, _ := reader.ReadString('\n')
		qtyInput = strings.TrimSpace(qtyInput)

		qty, err := strconv.Atoi(qtyInput)
		if err != nil || qty <= 0 {
			fmt.Println("Jumlah tidak valid.")
			continue
		}

		cart[kode] += qty
		fmt.Printf("✔ %d x %s ditambahkan\n", qty, item.Name)
	}

	// ================= HITUNG TOTAL =================
	var subtotal int64
	var lines []Line

	for code, qty := range cart {
		item, _ := findItemByCode(code)
		lineTotal := item.Price * int64(qty)
		subtotal += lineTotal

		lines = append(lines, Line{
			Item:      item,
			Qty:       qty,
			LineTotal: lineTotal,
		})
	}

	total := subtotal

	// ================= RINGKASAN =================
	fmt.Println("\n--- Ringkasan Pembelian ---")
	for _, l := range lines {
		fmt.Printf("%s x%d = %s\n",
			l.Item.Name,
			l.Qty,
			formatRupiah(l.LineTotal),
		)
	}
	fmt.Printf("TOTAL: %s\n", formatRupiah(total))

	// ================= PEMBAYARAN =================
	for {
		fmt.Print("Masukkan jumlah bayar: ")
		payInput, _ := reader.ReadString('\n')
		payInput = strings.TrimSpace(payInput)

		bayar, err := strconv.ParseInt(payInput, 10, 64)
		if err != nil {
			fmt.Println("Masukkan angka yang valid.")
			continue
		}

		// ✅ FITUR KEKURANGAN UANG
		if bayar < total {
			kurang := total - bayar
			fmt.Printf("Uang tidak cukup. Kurang %s\n", formatRupiah(kurang))
			continue
		}

		kembalian := bayar - total
		fmt.Println()
		printReceipt(lines, total, bayar, kembalian)
		break
	}
}

// ================== FUNGSI ==================

func printMenu() {
	fmt.Println("\n-- MENU --")
	for _, it := range menu {
		fmt.Printf("%d) %s - %s\n",
			it.Code,
			it.Name,
			formatRupiah(it.Price),
		)
	}
}

func findItemByCode(code int) (Item, bool) {
	for _, it := range menu {
		if it.Code == code {
			return it, true
		}
	}
	return Item{}, false
}

func formatRupiah(n int64) string {
	s := strconv.FormatInt(n, 10)
	var parts []string

	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)

	return "Rp " + strings.Join(parts, ".")
}

func printReceipt(lines []Line, total, bayar, kembalian int64) {
	fmt.Println("======================================")
	fmt.Println("       TOKO MAKANAN SEDERHANA")
	fmt.Println("======================================")
	fmt.Printf("Tanggal : %s\n",
		time.Now().Format("02-Jan-2006 15:04:05"))
	fmt.Println("--------------------------------------")
	fmt.Printf("%-20s %3s %10s\n", "Item", "Jumlah", "Total")
	fmt.Println("--------------------------------------")

	for _, l := range lines {
		name := l.Item.Name
		if len(name) > 20 {
			name = name[:20]
		}
		fmt.Printf("%-20s %3d %10s\n",
			name,
			l.Qty,
			formatRupiah(l.LineTotal),
		)
	}

	fmt.Println("--------------------------------------")
	fmt.Printf("%-26s %10s\n", "Total:", formatRupiah(total))
	fmt.Printf("%-26s %10s\n", "Bayar:", formatRupiah(bayar))
	fmt.Printf("%-26s %10s\n", "Kembalian:", formatRupiah(kembalian))
	fmt.Println("--------------------------------------")
	fmt.Println("Terima kasih telah berbelanja")
	fmt.Println("======================================")
}
