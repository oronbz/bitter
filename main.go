package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oronbz/bitter/internal/api"
	"github.com/oronbz/bitter/internal/config"
	"github.com/oronbz/bitter/internal/ui"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			runSetup()
			return
		case "help", "--help", "-h":
			printHelp()
			return
		case "version", "--version", "-v":
			fmt.Println("bitter " + version)
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "No Bitrise API token found.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Run `bitter setup` to configure your token interactively,")
		fmt.Fprintln(os.Stderr, "or run `bitter help` for all options.")
		os.Exit(1)
	}

	client := api.NewClient(cfg.Token)
	model := ui.NewModel(client)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`bitter — A terminal UI for Bitrise CI

Usage:
  bitter              Launch the TUI
  bitter setup        Configure your Bitrise API token
  bitter help         Show this help message
  bitter version      Show version

Environment:
  BITRISE_TOKEN       API token (overrides config file)

Config file:
  ~/.config/bitter/config.toml

Keybindings (inside TUI):
  j/k, ↑/↓            Navigate lists / scroll logs
  Enter                Select item
  Tab / Shift-Tab      Switch panel
  t                    Trigger new build
  a                    Abort running build
  r                    Refresh current view
  /                    Filter / search
  o                    Open build in browser
  c                    Copy build URL
  ?                    Show all keybindings
  q, Ctrl-C            Quit
`)
}

func runSetup() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot determine home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(home, ".config", "bitter")
	configPath := filepath.Join(configDir, "config.toml")

	// Check if already configured
	if cfg, err := config.Load(); err == nil && cfg.Token != "" {
		fmt.Println("You already have a token configured.")
		fmt.Print("Overwrite it? [y/N] ")
		answer := readLine()
		if !strings.HasPrefix(strings.ToLower(answer), "y") {
			fmt.Println("Setup cancelled.")
			return
		}
	}

	fmt.Println("bitter setup")
	fmt.Println("────────────")
	fmt.Println()
	fmt.Println("1. Open your Bitrise security settings to create a Personal Access Token:")
	fmt.Println()
	fmt.Println("   https://app.bitrise.io/me/account/security")
	fmt.Println()

	fmt.Print("Open in browser? [Y/n] ")
	answer := readLine()
	if answer == "" || strings.HasPrefix(strings.ToLower(answer), "y") {
		openBrowser("https://app.bitrise.io/me/account/security")
	}

	fmt.Println()
	fmt.Print("2. Paste your token here: ")
	token := readLine()
	token = strings.TrimSpace(token)

	if token == "" {
		fmt.Fprintln(os.Stderr, "No token provided. Setup cancelled.")
		os.Exit(1)
	}

	// Validate the token by making an API call
	fmt.Print("   Verifying token... ")
	client := api.NewClient(token)
	apps, _, err := client.ListApps("")
	if err != nil {
		fmt.Println("✗")
		fmt.Fprintf(os.Stderr, "\nToken validation failed: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please check your token and try again.")
		os.Exit(1)
	}
	fmt.Printf("✓ (%d apps found)\n", len(apps))

	// Save to config file
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	content := fmt.Sprintf("token = %q\n", token)
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("3. Token saved to %s\n", configPath)
	fmt.Println()
	fmt.Println("You're all set! Run `bitter` to launch.")
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Start()
}
