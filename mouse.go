package main

import (
    "fmt"
    "image"

    "gocv.io/x/gocv"
    "github.com/go-vgo/robotgo"

)

func main() {
    // Buka webcam
    webcam, err := gocv.OpenVideoCapture(0)
    if err != nil {
        fmt.Println("Error opening video capture device:", err)
        return
    }
    defer webcam.Close()

    // Buat window untuk menampilkan video
    window := gocv.NewWindow("Virtual Mouse")
    defer window.Close()

    // Inisialisasi MediaPipe Hands
    hands := gocv.NewMediaPipeHands()
    defer hands.Close()

    // Variabel untuk menyimpan riwayat koordinat
    var prevX, prevY int

    for {
        // Baca frame dari webcam
        img := webcam.Read()
        if img.Empty() {
            continue
        }

        // Deteksi tangan
        rects, _ := hands.Process(img)

        // Jika tangan terdeteksi
        if len(rects) > 0 {
            // Ambil landmark tangan pertama
            handLandmarks := hands.Hands(rects[0])

            // Dapatkan koordinat ujung jari telunjuk (landmark 8)
            indexFingerTip := handLandmarks[8]

            // Konversi koordinat landmark ke koordinat layar
            screenX, screenY := int(float64(indexFingerTip.X)*1.5), int(float64(indexFingerTip.Y)*1.5) // Sesuaikan faktor skala jika perlu

            // Smoothing dengan rata-rata bergerak sederhana
            screenX = (screenX + prevX) / 2
            screenY = (screenY + prevY) / 2

            // Simpan koordinat saat ini untuk iterasi berikutnya
            prevX, prevY = screenX, screenY

            // Gerakkan kursor mouse
            robotgo.MoveMouse(screenX, screenY)

            // Contoh logika klik (jika jari telunjuk dan ibu jari berdekatan)
            thumbTip := handLandmarks[4]
            distance := image.Pt(indexFingerTip.X, indexFingerTip.Y).Sub(image.Pt(thumbTip.X, thumbTip.Y)).Norm()
            if distance < 30 { // Atur ambang jarak sesuai kebutuhan
                robotgo.Click("left")
            }
        }

        // Tampilkan frame dengan landmark
        hands.DrawHands(img, rects)
        window.IMShow(img)

        // Keluar jika tombol 'q' ditekan
        if window.WaitKey(1) == 'q' {
            break
        }
    }
}
