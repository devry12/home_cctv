package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Membuat channel untuk mengirim sinyal ketika rekaman selesai
	done := make(chan struct{})

	// Membuat slice untuk menyimpan konfigurasi kamera
	cctvs := []struct {
		Name     string
		Url      string
		Username string
		Password string
		Duration string
	}{
		{"cctv1", "192.168.2.89", "admin", "devry123", "60"},
		{"cctv2", "192.168.2.15", "admin", "devry123", "60"},
	}

	// Memulai goroutine untuk merekam video secara berulang
	for _, cctv := range cctvs {
		go func(cctv struct {
			Name     string
			Url      string
			Username string
			Password string
			Duration string
		}) {
			for {
				err := recordVideo(cctv)
				if err != nil {
					fmt.Printf("Error recording video for %s: %v\n", cctv.Name, err)
					return
				}
				// Mengirim sinyal ke channel 'done' ketika rekaman selesai
				done <- struct{}{}
			}
		}(cctv)
	}

	// Menunggu sinyal dari channel 'done' (akan selalu menerima sinyal)
	for {
		<-done
	}
}

func recordVideo(cctv struct {
	Name     string
	Url      string
	Username string
	Password string
	Duration string
}) error {
	// Membuat nama folder dengan format tanggal saat ini (YYYY-MM-DD)
	folderName := time.Now().Format("2006-01-02")

	// Buat direktori nama cctv jika belum ada
	err := os.MkdirAll(cctv.Name+"/"+folderName, 0755)
	if err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	// Buat nama file dengan format waktu saat ini (HH-MM-SS).mp4
	fileName := time.Now().Format("15-04") + ".mp4"

	// Path lengkap ke file output
	outputPath := fmt.Sprintf("%s/%s/%s", cctv.Name, folderName, fileName)

	// Command and arguments are separate elements in a slice
	cmd := exec.Command("ffmpeg", "-rtsp_transport", "udp", "-i", fmt.Sprintf("rtsp://%s:%s@%s/onvif1", cctv.Username, cctv.Password, cctv.Url), "-t", cctv.Duration, outputPath)

	// Membuat pipe untuk menangkap output STDERR
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Jalankan perintah
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting ffmpeg command: %v", err)
	}

	// Buat goroutine untuk menampilkan output STDERR
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	// Tunggu perintah selesai
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Printf("Output saved to: %s\n", outputPath)
	return nil
}
