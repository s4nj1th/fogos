package cli

import (
	"fmt"
	"os"
	"strings"

	"fogos/internal/blocker"
	"fogos/pkg/validator"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

type App struct {
	blocker *blocker.Blocker
}

func NewApp() *App {
	return &App{
		blocker: blocker.New(),
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		return fmt.Errorf("insufficient arguments")
	}

	command := args[1]
	var website string
	if len(args) > 2 {
		website = strings.TrimSpace(args[2])
	}

	switch command {
	case "b":
		command = "block"
	case "ub":
		command = "unblock"
	case "s":
		command = "status"
	case "l":
		command = "list"
	}

	if command != "list" {
		if command != "status" && command != "unblock" && command != "block" {
			fmt.Printf("%sUnknown command:%s %s\n", ColorRed, ColorReset, command)
			a.printUsage()
			return fmt.Errorf("unknown command: %s", command)
		}
		if website == "" {
			a.printUsage()
			return fmt.Errorf("missing website argument")
		}
		if command != "status" {
			if err := validator.ValidateWebsite(website); err != nil {
				fmt.Printf("%sError:%s %v\n", ColorRed, ColorReset, err)
				return err
			}
		}
	}

	if (command == "block" || command == "unblock") && os.Geteuid() != 0 {
		fmt.Printf("%sError:%s This program must be run as root/administrator\n", ColorRed, ColorReset)
		fmt.Printf("Try: sudo fogos %s %s\n", command, website)
		return fmt.Errorf("insufficient privileges")
	}

	switch command {
	case "block":
		return a.handleBlock(website)
	case "unblock":
		return a.handleUnblock(website)
	case "status":
		return a.handleStatus(website)
	case "list":
		return a.handleList()
	default:
		fmt.Printf("%sUnknown command:%s %s\n", ColorRed, ColorReset, command)
		a.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) handleBlock(website string) error {
	if err := a.blocker.Block(website); err != nil {
		fmt.Printf("%sError blocking %s:%s %v\n", ColorRed, website, ColorReset, err)
		return err
	}
	fmt.Printf("%sSuccessfully blocked %s%s\n", ColorGreen, website, ColorReset)
	return nil
}

func (a *App) handleUnblock(website string) error {
	if err := a.blocker.Unblock(website); err != nil {
		fmt.Printf("%sError unblocking %s:%s %v\n", ColorRed, website, ColorReset, err)
		return err
	}
	fmt.Printf("%sSuccessfully unblocked %s%s\n", ColorGreen, website, ColorReset)
	return nil
}

func (a *App) handleStatus(website string) error {
	blocked, err := a.blocker.IsBlocked(website)
	if err != nil {
		fmt.Printf("%sError checking status of %s:%s %v\n", ColorRed, website, ColorReset, err)
		return err
	}
	if blocked {
		fmt.Printf("%s%s is currently blocked%s\n", ColorRed, website, ColorReset)
	} else {
		fmt.Printf("%s%s is currently accessible%s\n", ColorGreen, website, ColorReset)
	}
	return nil
}

func (a *App) handleList() error {
	websites, err := a.blocker.ListBlocked()
	if err != nil {
		fmt.Printf("%sError listing blocked websites:%s %v\n", ColorRed, ColorReset, err)
		return err
	}
	if len(websites) == 0 {
		fmt.Printf("%sNo websites are currently blocked%s\n", ColorGreen, ColorReset)
		return nil
	}
	fmt.Printf("%sBlocked websites (%d):%s\n", ColorYellow, len(websites), ColorReset)
	for _, website := range websites {
		fmt.Printf("   %s%s%s\n", ColorRed, website, ColorReset)
	}
	return nil
}

func (a *App) printUsage() {
	fmt.Println("fogos - Website blocker using /etc/hosts")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  fogos block [website]     (alias: b) Block a website")
	fmt.Println("  fogos unblock [website]   (alias: ub) Unblock a website")
	fmt.Println("  fogos status [website]    (alias: s) Check if website is blocked")
	fmt.Println("  fogos list                (alias: l) List all blocked websites")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  sudo fogos block youtube.com")
	fmt.Println("  sudo fogos unblock youtube.com")
	fmt.Println("  fogos status youtube.com")
	fmt.Println("  fogos list")
}
