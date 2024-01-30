package Models

import (
	"GoIris/Models/Database"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/disk"
	"golang.org/x/text/message"
)

func ClearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func SetupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		exit()
	}()
}

func ConvertBytesToPretty(bytes int64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))

	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes < KB:
		return p.Sprintf("%d B", bytes)
	case bytes < MB:
		return p.Sprintf("%.2f KB", float64(bytes)/KB)
	case bytes < GB:
		return p.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes < TB:
		return p.Sprintf("%.2f GB", float64(bytes)/GB)
	default:
		return p.Sprintf("%.2f TB", float64(bytes)/TB)
	}
}

func UpdateInsertConsoleTitle(ctx context.Context) {
	p := message.NewPrinter(message.MatchLanguage("en"))
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Exit the function when context is cancelled
		case <-ticker.C:
			Database.Mutex.Lock()
			totalInserts := Database.TotalInserts
			Database.Mutex.Unlock()

			elapsedTime := time.Since(Database.StartTime)
			var insertsPerSecond float64
			if elapsedTime.Seconds() > 0 {
				insertsPerSecond = float64(totalInserts) / elapsedTime.Seconds()
			}

			formattedElapsedTime := elapsedTime.Truncate(time.Second).String()
			formattedTotalInserts := p.Sprintf("%d", totalInserts)
			formattedInsertsPerSecond := p.Sprintf("%.2f", insertsPerSecond)

			title := fmt.Sprintf("GoIris ~ Total Inserts: %s - Inserts/sec: %s - Elapsed Time: %s", formattedTotalInserts, formattedInsertsPerSecond, formattedElapsedTime)
			fmt.Printf("\033]0;%s\007", title)
		}
	}
}

func CheckDiskSpace() uint64 {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diskStat, err := disk.Usage(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return diskStat.Free
}

func exit() {
	fmt.Println("Performing cleanup...")
	Database.SetPragmasForExit(DB)
	fmt.Println("Exiting...")
	os.Exit(0)
}
