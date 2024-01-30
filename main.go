package main

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	// "GoIris/Models/Hashes"
	"GoIris/Models/Modules"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/sqweek/dialog"
	_ "modernc.org/sqlite"

)

var TotalInserts int
var StartTime time.Time

func main() {
	Models.SetupSignalHandler()
	reader := bufio.NewReader(os.Stdin)
	for {
		Models.ClearConsole()
		fmt.Printf("\033]0;%s\007", "GoIris ~ GoLang Rainbow Table")
		choice := promptUser(reader)
		switch choice {
		case "1":
			ctx, cancel := context.WithCancel(context.Background())
			Database.SetPragmasForInsert(Models.DB)
			go Models.UpdateInsertConsoleTitle(ctx)

			err := Modules.InsertFile()
			if err != nil {
				fmt.Println("Error opening file dialog:", err)
			}

			cancel()
		case "2":
			Database.SetPragmasForRead(Models.DB)
			Modules.DehashFile()
		case "3":
			Database.SetPragmasForRead(Models.DB)
			Modules.LookupHash(reader)
		case "4":
			fmt.Println("Exiting...")
			return
		}
	}
}

func promptUser(reader *bufio.Reader) string {
	fmt.Println("What would you like to do?")
	fmt.Println("1. Insert a file")
	fmt.Println("2. Dehash a File")
	fmt.Println("3. Lookup a hash")
	fmt.Println("4. Exit")

	fmt.Print("Enter choice: ")
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}
