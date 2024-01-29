package Models

import (
	"GoIris/Models/Database"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

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

func exit() {
	fmt.Println("Performing cleanup...")
	Database.SetPragmasForExit(DB)
	fmt.Println("Exiting...")
	os.Exit(0)
}

func SetupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		exit()
	}()
}

func UpdateInsertConsoleTitle() {
	p := message.NewPrinter(message.MatchLanguage("en"))

	for range time.Tick(2 * time.Second) {
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
