package localfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FileItem struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
}

type Model struct {
	Files       []FileItem
	SelectedIdx int
	Dir         string
	Width       int
	Height      int
	ErrMsg      string
}

type ItemsLoadedMsg struct {
	Files []FileItem
	Err   string
}

type PlayFileMsg struct {
	Path string
}

func NewModel(dir string) Model {
	return Model{Dir: dir}
}

func (m Model) Load() tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(m.Dir)
		if err != nil {
			return ItemsLoadedMsg{Err: err.Error()}
		}
		var files []FileItem
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(e.Name()))
			switch ext {
			case ".mp4", ".mkv", ".webm", ".mp3", ".m4a", ".ogg", ".flac":
				info, _ := e.Info()
				files = append(files, FileItem{
					Name:    e.Name(),
					Path:    filepath.Join(m.Dir, e.Name()),
					Size:    info.Size(),
					ModTime: info.ModTime(),
				})
			}
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime.After(files[j].ModTime)
		})
		return ItemsLoadedMsg{Files: files}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) HandleResize(w, h int) Model {
	m.Width = w
	m.Height = h
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ItemsLoadedMsg:
		if msg.Err != "" {
			m.ErrMsg = msg.Err
		} else {
			m.Files = msg.Files
			m.SelectedIdx = 0
		}
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			if m.SelectedIdx < len(m.Files)-1 {
				m.SelectedIdx++
			}
		case "k", "up":
			if m.SelectedIdx > 0 {
				m.SelectedIdx--
			}
		case "enter":
			if m.SelectedIdx < len(m.Files) {
				path := m.Files[m.SelectedIdx].Path
				return m, func() tea.Msg { return PlayFileMsg{Path: path} }
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Local Files (%s)", m.Dir)))
	s.WriteRune('\n')

	if m.ErrMsg != "" {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render("⚠ " + m.ErrMsg))
		return s.String()
	}
	if len(m.Files) == 0 {
		s.WriteString("  (no media files found)")
		return s.String()
	}
	maxRows := m.Height - 3
	if m.Height == 0 {
		maxRows = 15
	}
	if maxRows < 1 {
		maxRows = 1
	}
	start := 0
	if m.SelectedIdx >= maxRows {
		start = m.SelectedIdx - maxRows + 1
	}
	end := start + maxRows
	if end > len(m.Files) {
		end = len(m.Files)
	}
		for i := start; i < end; i++ {
		f := m.Files[i]
		prefix := "  "
		if i == m.SelectedIdx {
			prefix = "❯ "
		}
		sizeStr := fmt.Sprintf("%.1f MB", float64(f.Size)/(1024*1024))
		s.WriteString(fmt.Sprintf("%s%s  %s  %s\n", prefix, f.Name, sizeStr, f.ModTime.Format("Jan 2 15:04")))
	}
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("  j/k navigate · enter play · esc back"))
	return s.String()
}   
