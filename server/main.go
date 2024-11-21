package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// Struct untuk merepresentasikan pengguna
type User struct {
	Name  string
	Saldo int
}

var (
	users    = make(map[string]*User)         // Peta untuk menyimpan data pengguna
	clients  = make(map[*websocket.Conn]bool) // Peta untuk menyimpan koneksi WebSocket
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Mengizinkan semua asal untuk contoh ini
		},
	}
)

// Fungsi untuk mengirim pesan ke semua klien WebSocket
func broadcastMessage(message string) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(clients, client) // Menghapus klien jika terjadi kesalahan
		}
	}
}

// Handler untuk koneksi WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Meningkatkan koneksi HTTP ke WebSocket
	if err != nil {
		fmt.Println("Error saat meng-upgrade ke WebSocket:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true // Menyimpan koneksi klien baru

	for {
		_, _, err := conn.ReadMessage() // Menunggu pesan dari klien
		if err != nil {
			delete(clients, conn) // Menghapus klien jika terjadi kesalahan saat membaca
			break
		}
	}
}

// Fungsi untuk menangani koneksi UDP
func prosesClientUDP(conn *net.UDPConn) {
	buffer := make([]byte, 1024) // Buffer untuk menerima data
	for {
		n, _, err := conn.ReadFromUDP(buffer) // Membaca data dari koneksi UDP
		if err != nil {
			fmt.Println("Error saat menerima data UDP:", err)
			continue
		}

		message := string(buffer[:n])       // Mengonversi data yang diterima ke string
		commands := strings.Fields(message) // Memisahkan pesan menjadi array kata
		if len(commands) < 2 {
			fmt.Println("Format pesan tidak valid.")
			continue
		}

		switch commands[0] {
		case "REGISTER":
			name := commands[1]

			if _, exists := users[name]; !exists {
				users[name] = &User{Name: name, Saldo: 0} // Menambahkan pengguna baru
				fmt.Printf("Pengguna %s terdaftar dengan saldo 0.\n", name)
			} else {
				fmt.Printf("Pengguna %s sudah terdaftar.\n", name)
			}

		case "DONATE":
			sender := commands[1]
			jumlahDonasi, err := strconv.Atoi(commands[2]) // Mengonversi jumlah donasi ke int
			if err != nil {
				fmt.Println("Jumlah donasi tidak valid.")
				continue
			}

			donationMessage := ""
			if len(commands) > 3 {
				donationMessage = strings.Join(commands[3:], " ") // Menggabungkan pesan donasi
			}

			if user, exists := users[sender]; exists {
				if user.Saldo >= jumlahDonasi {
					user.Saldo -= jumlahDonasi // Mengurangi saldo pengguna
					fmt.Printf("%s mendonasikan %d. Pesan: \"%s\". Saldo sekarang: %d\n", sender, jumlahDonasi, donationMessage, user.Saldo)
					broadcastMessage(fmt.Sprintf("%s mendonasikan %d. Pesan: \"%s\"", sender, jumlahDonasi, donationMessage))
				} else {
					fmt.Printf("Saldo %s tidak mencukupi untuk donasi.\n", sender)
				}
			} else {
				fmt.Printf("Pengguna %s tidak ditemukan.\n", sender)
			}
		}
	}
}

// Fungsi untuk menangani koneksi TCP
func prosesTopUpTCP(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer) // Membaca data dari koneksi TCP
	if err != nil {
		fmt.Println("Error saat membaca data TCP:", err)
		return
	}

	message := string(buffer[:n])       // Mengonversi data yang diterima ke string
	commands := strings.Fields(message) // Memisahkan pesan menjadi array kata
	if len(commands) < 2 {
		fmt.Println("Format TOPUP tidak valid.")
		return
	}

	name := commands[0]
	jumlahDonasi, err := strconv.Atoi(commands[1]) // Mengonversi jumlah top up ke int
	if err != nil {
		fmt.Println("Jumlah top up tidak valid.")
		return
	}

	if user, exists := users[name]; exists {
		user.Saldo += jumlahDonasi // Menambah saldo pengguna
		fmt.Printf("Saldo %s ditambahkan %d. Saldo sekarang: %d.\n", name, jumlahDonasi, user.Saldo)
	} else {
		fmt.Printf("Pengguna %s tidak ditemukan.\n", name)
	}
}

func main() {
	// Server WebSocket untuk klien HTML
	http.HandleFunc("/ws", handleWebSocket)
	go func() {
		if err := http.ListenAndServe(":8082", nil); err != nil {
			fmt.Println("Error saat menjalankan server WebSocket:", err)
		}
	}()
	fmt.Println("Server berjalan (WebSocket: 8082, UDP: 8080 untuk registrasi & donasi, TCP: 8081 untuk top up)")

	// Setup UDP
	udpAddr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("127.0.0.1"),
	}
	udpConn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		fmt.Println("Error saat membuat koneksi UDP:", err)
		return
	}
	defer udpConn.Close()
	go prosesClientUDP(udpConn)

	// Setup TCP
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error saat membuat koneksi TCP:", err)
		return
	}
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.Accept() // Menerima koneksi TCP baru
		if err != nil {
			fmt.Println("Error saat menerima koneksi TCP:", err)
			continue
		}
		go prosesTopUpTCP(tcpConn) // Menangani koneksi dalam goroutine baru
	}
}
