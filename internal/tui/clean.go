package tui

import (
	"sort"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/iyushkarki/csys/internal/cleaners"
	"github.com/iyushkarki/csys/internal/system"
)

const cleanWorkers = 3

type phase int

const (
	selecting phase = iota
	cleaning
	done
)

type itemStatus int

const (
	statusIdle itemStatus = iota
	statusWorking
	statusDone
	statusFailed
)

type item struct {
	target   *cleaners.Target
	selected bool
	status   itemStatus
	err      error
}

type targetFoundMsg struct{ target *cleaners.Target }
type scanDoneMsg struct{}
type cleanedMsg struct {
	it  *item
	err error
}

type Model struct {
	phase    phase
	opts     cleaners.Options
	spin     spinner.Model
	items    []*item
	cursor   int
	scanCh   chan *cleaners.Target
	scanning bool
	width    int
	height   int

	jobs       chan *item
	pending    int
	freed      uint64
	trashed    uint64
	diskBefore *system.DiskInfo
	diskAfter  *system.DiskInfo
}

func Run(produce func() []*cleaners.Target, opts cleaners.Options) error {
	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(accent)

	m := &Model{
		opts:     opts,
		spin:     sp,
		scanCh:   make(chan *cleaners.Target, 256),
		scanning: true,
		width:    80,
		height:   30,
	}
	go func() { cleaners.DetectStream(produce(), m.scanCh) }()

	_, err := tea.NewProgram(m).Run()
	return err
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, waitForTarget(m.scanCh))
}

func waitForTarget(ch chan *cleaners.Target) tea.Cmd {
	return func() tea.Msg {
		t, ok := <-ch
		if !ok {
			return scanDoneMsg{}
		}
		return targetFoundMsg{t}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd

	case targetFoundMsg:
		m.items = append(m.items, &item{target: msg.target, selected: msg.target.Preselect})
		m.sortItems()
		return m, waitForTarget(m.scanCh)

	case scanDoneMsg:
		m.scanning = false
		if len(m.items) == 0 {
			m.phase = done
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case cleanedMsg:
		msg.it.err = msg.err
		if msg.err != nil {
			msg.it.status = statusFailed
		} else {
			msg.it.status = statusDone
			if msg.it.target.UsesTrash(m.opts) {
				m.trashed += msg.it.target.Size
			} else {
				m.freed += msg.it.target.Size
			}
		}
		m.pending--
		if m.pending == 0 {
			m.finish()
			return m, nil
		}
		return m, m.workCmd()
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.phase == cleaning {
		return m, nil
	}
	if key == "ctrl+c" {
		return m, tea.Quit
	}
	if m.phase == done {
		if key == "q" || key == "enter" || key == "esc" {
			return m, tea.Quit
		}
		return m, nil
	}

	switch key {
	case "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case " ":
		if len(m.items) > 0 {
			m.items[m.cursor].selected = !m.items[m.cursor].selected
		}
	case "a":
		for _, it := range m.items {
			if it.target.Tier == cleaners.Safe {
				it.selected = true
			}
		}
	case "u":
		for _, it := range m.items {
			it.selected = false
		}
	case "enter":
		return m, m.startCleaning()
	}
	return m, nil
}

func (m *Model) startCleaning() tea.Cmd {
	var selected []*item
	for _, it := range m.items {
		if it.selected {
			selected = append(selected, it)
		}
	}
	if len(selected) == 0 {
		return nil
	}

	m.phase = cleaning
	m.pending = len(selected)
	m.diskBefore, _ = system.GetDiskInfo()
	m.jobs = make(chan *item, len(selected))
	for _, it := range selected {
		it.status = statusWorking
		m.jobs <- it
	}
	close(m.jobs)

	workers := cleanWorkers
	if workers > len(selected) {
		workers = len(selected)
	}
	cmds := []tea.Cmd{m.spin.Tick}
	for w := 0; w < workers; w++ {
		cmds = append(cmds, m.workCmd())
	}
	return tea.Batch(cmds...)
}

func (m *Model) workCmd() tea.Cmd {
	return func() tea.Msg {
		it, ok := <-m.jobs
		if !ok {
			return nil
		}
		return cleanedMsg{it: it, err: it.target.Clean(m.opts)}
	}
}

func (m *Model) finish() {
	m.phase = done
	m.diskAfter, _ = system.GetDiskInfo()

	var entries []cleaners.HistoryEntry
	for _, it := range m.items {
		if it.status == statusDone {
			entries = append(entries, cleaners.NewHistoryEntry(it.target, m.opts))
		}
	}
	cleaners.AppendHistory(entries)
}

func (m *Model) sortItems() {
	current := m.currentItem()
	sort.SliceStable(m.items, func(i, j int) bool {
		return cleaners.Less(m.items[i].target, m.items[j].target)
	})
	if current != nil {
		for i, it := range m.items {
			if it == current {
				m.cursor = i
				break
			}
		}
	}
}

func (m *Model) currentItem() *item {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return m.items[m.cursor]
	}
	return nil
}

func (m *Model) selectedSize() uint64 {
	var total uint64
	for _, it := range m.items {
		if it.selected {
			total += it.target.Size
		}
	}
	return total
}

func (m *Model) foundSize() uint64 {
	var total uint64
	for _, it := range m.items {
		total += it.target.Size
	}
	return total
}
