package main

import (
	"bufio"   // Untuk membaca input dari pengguna
	"fmt"     // Untuk mencetak ke layar
	"net"     // Untuk menggunakan fungsi jaringan
	"os"      // Untuk interaksi dengan sistem operasi, seperti membaca input
	"strings" // Untuk manipulasi string
)

func main() {
	// Membuat alamat server UDP
	serverAddr := net.UDPAddr{
		Port: 8080,                     // Port server yang digunakan
		IP:   net.ParseIP("127.0.0.1"), // IP server (localhost)
	}

	// Membuka koneksi UDP ke server
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		fmt.Println("Error:", err) // Menampilkan error jika koneksi gagal
		return
	}
	defer conn.Close() // Menutup koneksi ketika program selesai

	reader := bufio.NewReader(os.Stdin) // Membuat reader untuk membaca input dari terminal
	fmt.Print("Masukkan nama Anda: ")
	name, _ := reader.ReadString('\n') // Membaca input nama dari pengguna
	name = strings.TrimSpace(name)     // Menghapus spasi di awal/akhir string

	// Mengirim pesan registrasi ke server
	registerMessage := "REGISTER " + name        // Format pesan registrasi
	_, err = conn.Write([]byte(registerMessage)) // Mengirim data ke server
	if err != nil {
		fmt.Println("Error saat mengirim data:", err) // Menampilkan error jika pengiriman gagal
		return
	}
	fmt.Println("Pendaftaran berhasil!") // Memberitahu pengguna jika pendaftaran berhasil

	// Loop utama untuk pilihan pengguna
	for {
		fmt.Println("\nPilih opsi:")
		fmt.Println("1. Donasi")
		fmt.Println("2. Keluar")
		fmt.Print("Masukkan pilihan Anda: ")
		choice, _ := reader.ReadString('\n') // Membaca pilihan pengguna
		choice = strings.TrimSpace(choice)   // Menghapus spasi di awal/akhir string

		switch choice { // Memeriksa input pilihan pengguna
		case "1":
			fmt.Print("Masukkan jumlah donasi: ")
			amount, _ := reader.ReadString('\n') // Membaca input jumlah donasi
			amount = strings.TrimSpace(amount)   // Menghapus spasi di awal/akhir string

			fmt.Print("Masukkan pesan untuk donasi Anda: ")
			message, _ := reader.ReadString('\n') // Membaca input pesan donasi
			message = strings.TrimSpace(message)  // Menghapus spasi di awal/akhir string

			// Mengirim pesan donasi ke server dalam format "DONATE <nama> <jumlah> <pesan>"
			donateMessage := fmt.Sprintf("DONATE %s %s %s", name, amount, message)
			_, err = conn.Write([]byte(donateMessage)) // Mengirim data ke server
			if err != nil {
				fmt.Println("Error saat mengirim donasi:", err) // Menampilkan error jika pengiriman gagal
			} else {
				fmt.Println("Donasi berhasil dikirim!") // Memberitahu pengguna jika donasi berhasil dikirim
			}

		case "2":
			fmt.Println("Keluar dari program.") // Memberitahu pengguna bahwa program akan keluar
			return                              // Menghentikan program

		default:
			fmt.Println("Pilihan tidak valid.") // Memberitahu pengguna jika input tidak valid
		}
	}
}
