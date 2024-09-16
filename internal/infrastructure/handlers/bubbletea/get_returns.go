package bubbletea

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"homework/internal/abstractions"
	"strconv"
)

type getReturnsModel struct {
	useCase abstractions.IPVZOrderUseCase

	table table.Model

	page     int
	pageSize int

	changed bool
}

func newGetReturnsModel(useCase abstractions.IPVZOrderUseCase, pageSize int) *getReturnsModel {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Recipient ID", Width: 15},
		{Title: "Order ID", Width: 10},
		{Title: "ReceivedAt", Width: 20},
		{Title: "StorageTime", Width: 15},
		{Title: "IssuedAt", Width: 20},
		{Title: "ReturnedAt", Width: 20},
		{Title: "Weight", Width: 10},
		{Title: "Cost", Width: 10},
		{Title: "Packaging", Width: 10},
		{Title: "AdditionalFilm", Width: 15},
	}
	dataTable := table.New(
		table.WithColumns(columns),
		table.WithHeight(pageSize),
		table.WithFocused(false),
	)

	return &getReturnsModel{
		useCase:  useCase,
		table:    dataTable,
		pageSize: pageSize,
		changed:  true,
	}
}

func (m *getReturnsModel) Init() tea.Cmd {
	return nil
}

func (m *getReturnsModel) handleKeyboard(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return tea.Quit
	case tea.KeyDown:
		m.page++
		m.changed = true
	case tea.KeyUp:
		m.page--
		if m.page < 0 {
			m.page = 0
		}
		m.changed = true
	default:
	}

	return nil
}

func (m *getReturnsModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	switch tea.MouseEvent(msg).Button {
	case tea.MouseButtonWheelDown:
		m.page++
		m.changed = true
	case tea.MouseButtonWheelUp:
		m.page--
		if m.page < 0 {
			m.page = 0
		}
		m.changed = true
	default:
	}

	return nil
}

func (m *getReturnsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = append(cmds, m.handleKeyboard(msg))
	case tea.MouseMsg:
		cmds = append(cmds, m.handleMouse(msg))
	default:
	}

	return m, tea.Batch(cmds...)
}

func (m *getReturnsModel) View() string {
	if m.changed {
		orders, err := m.useCase.GetReturns(abstractions.WithPage(m.page), abstractions.WithPageSize(m.pageSize))
		if err != nil {
			return err.Error()
		}
		rows := make([]table.Row, len(orders))
		for i, order := range orders {
			rows[i] = table.Row{
				order.OrderID,
				order.PVZID,
				order.RecipientID,
				order.ReceivedAt.Format("2006-01-02 15:04:05"),
				order.StorageTime.String(),
				order.IssuedAt.Format("2006-01-02 15:04:05"),
				order.ReturnedAt.Format("2006-01-02 15:04:05"),
				strconv.Itoa(order.Weight),
				strconv.Itoa(order.Cost),
				order.Packaging.String(),
				strconv.FormatBool(order.AdditionalFilm),
			}
		}
		m.table.SetRows(rows)
		m.changed = false
	}

	return m.table.View() + "\n" + "Page: " + strconv.Itoa(m.page+1) + "\n" + m.table.HelpView()
}
