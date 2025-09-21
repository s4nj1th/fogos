package hosts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

import "runtime"

const (
	FogosComment = "# Blocked by fogos"
)

func getDefaultHostsFile() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

type Entry struct {
	IP      string
	Domain  string
	Comment string
}

func (e Entry) String() string {
	if e.Comment != "" {
		return fmt.Sprintf("%s %s %s", e.IP, e.Domain, e.Comment)
	}
	return fmt.Sprintf("%s %s", e.IP, e.Domain)
}

type Manager struct {
	filePath string
}


func New() *Manager {
	return &Manager{
		filePath: getDefaultHostsFile(),
	}
}

func NewWithPath(path string) *Manager {
	return &Manager{
		filePath: path,
	}
}

func (m *Manager) AddEntries(entries []Entry) error {
	file, err := os.OpenFile(m.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open hosts file: %v", err)
	}
	defer file.Close()

	for _, entry := range entries {
		if _, err := file.WriteString(entry.String() + "\n"); err != nil {
			return fmt.Errorf("failed to write to hosts file: %v", err)
		}
	}

	return nil
}

func (m *Manager) RemoveEntriesByDomain(domain string) error {
	lines, err := m.readLines()
	if err != nil {
		return err
	}

	var newLines []string
	removed := false

	for _, line := range lines {
		if m.isFogosEntry(line) && m.containsDomain(line, domain) {
			removed = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !removed {
		return fmt.Errorf("no entries found for domain: %s", domain)
	}

	return m.writeLines(newLines)
}

func (m *Manager) GetFogosEntries() ([]Entry, error) {
	lines, err := m.readLines()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, line := range lines {
		if m.isFogosEntry(line) {
			if entry, err := m.parseLine(line); err == nil {
				entries = append(entries, entry)
			}
		}
	}

	return entries, nil
}

func (m *Manager) readLines() ([]string, error) {
	file, err := os.Open(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open hosts file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %v", err)
	}

	return lines, nil
}

func (m *Manager) writeLines(lines []string) error {
	file, err := os.OpenFile(m.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open hosts file for writing: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to hosts file: %v", err)
		}
	}

	return writer.Flush()
}

func (m *Manager) isFogosEntry(line string) bool {
	return strings.Contains(line, FogosComment)
}

func (m *Manager) containsDomain(line, domain string) bool {
	return strings.Contains(line, " "+domain+" ") ||
		strings.Contains(line, " www."+domain+" ") ||
		strings.HasSuffix(line, " "+domain+" "+FogosComment) ||
		strings.HasSuffix(line, " www."+domain+" "+FogosComment)
}

func (m *Manager) parseLine(line string) (Entry, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Entry{}, fmt.Errorf("invalid hosts entry: %s", line)
	}

	entry := Entry{
		IP:     parts[0],
		Domain: parts[1],
	}

	if commentIndex := strings.Index(line, "#"); commentIndex >= 0 {
		entry.Comment = strings.TrimSpace(line[commentIndex:])
	}

	return entry, nil
}
