package main

import (
	"bufio"   // Paket untuk membaca input dari pengguna
	"fmt"     // Paket untuk mencetak output ke layar
	"net"     // Paket untuk melakukan komunikasi jaringan
	"os"      // Paket untuk interaksi dengan sistem operasi
	"strings" // Paket untuk manipulasi string, seperti trim
)

func main() {
	// Membuka koneksi TCP ke server dengan alamat IP 127.0.0.1 dan port 8081
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("Error:", err) // Menampilkan pesan kesalahan jika koneksi gagal
		return
	}
	defer conn.Close() // Menutup koneksi ketika fungsi `main` selesai dieksekusi

	// Membuat pembaca (reader) untuk membaca input dari terminal
	reader := bufio.NewReader(os.Stdin)

	// Meminta pengguna memasukkan nama pengguna yang ingin di-top up
	fmt.Print("Masukkan nama pengguna yang ingin di-top up: ")
	name, _ := reader.ReadString('\n') // Membaca input pengguna sampai newline
	name = strings.TrimSpace(name)     // Menghapus spasi atau newline di awal/akhir input

	// Meminta pengguna memasukkan jumlah top up
	fmt.Print("Masukkan jumlah top up: ")
	amount, _ := reader.ReadString('\n') // Membaca input jumlah sampai newline
	amount = strings.TrimSpace(amount)   // Menghapus spasi atau newline di awal/akhir input

	// Membuat pesan top up dengan format "nama jumlah"
	pesanTopUp := name + " " + amount

	// Mengirim pesan top up ke server
	_, err = conn.Write([]byte(pesanTopUp))
	if err != nil {
		fmt.Println("Error saat mengirim top up:", err) // Menampilkan pesan kesalahan jika pengiriman gagal
	} else {
		fmt.Println("Top up berhasil dikirim!") // Menampilkan pesan sukses jika pengiriman berhasil
	}
}
