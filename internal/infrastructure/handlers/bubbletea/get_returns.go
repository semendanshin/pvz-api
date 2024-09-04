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

func (m *getReturnsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
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
	case tea.MouseMsg:
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
	default:
	}

	return m, nil
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
			}
		}
		m.table.SetRows(rows)
		m.changed = false
	}

	return m.table.View() + "\n" + "Page: " + strconv.Itoa(m.page+1) + "\n" + m.table.HelpView()
}
