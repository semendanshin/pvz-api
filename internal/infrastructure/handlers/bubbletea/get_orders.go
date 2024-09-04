package bubbletea

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"homework/internal/abstractions"
	"strconv"
)

var _ tea.Model = &getOrdersModel{}

// getOrdersModel is a model for getting orders
type getOrdersModel struct {
	useCase abstractions.IPVZOrderUseCase

	userIDInput textinput.Model
	table       table.Model

	userID string

	changed  bool
	page     int
	pageSize int
}

// newGetOrdersModel creates a new getOrdersModel
func newGetOrdersModel(useCase abstractions.IPVZOrderUseCase, pageSize int) *getOrdersModel {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "PVZ ID", Width: 10},
		{Title: "Recipient ID", Width: 15},
		{Title: "ReceivedAt", Width: 20},
		{Title: "StorageTime", Width: 15},
		{Title: "IssuedAt", Width: 20},
		{Title: "ReturnedAt", Width: 20},
	}
	dataTable := table.New(
		table.WithColumns(columns),
		table.WithHeight(pageSize),
		table.WithFocused(true),
	)

	input := textinput.New()
	input.Placeholder = "Enter user ID"
	input.Prompt = "User ID: "
	input.Focus()

	return &getOrdersModel{
		useCase:     useCase,
		pageSize:    pageSize,
		table:       dataTable,
		userIDInput: input,
	}
}

// Init initializes the model
func (m *getOrdersModel) Init() tea.Cmd {
	return nil
}

// Update updates the model
func (m *getOrdersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if !m.userIDInput.Focused() {
				m.userIDInput.Focus()
			} else {
				m.userIDInput.SetValue("")
				return m, tea.Quit
			}
		case tea.KeyDown:
			m.page++
			m.changed = true
		case tea.KeyUp:
			m.page--
			if m.page < 0 {
				m.page = 0
			}
			m.changed = true
		case tea.KeyEnter:
			if m.userIDInput.Focused() {
				m.userID = m.userIDInput.Value()
				m.changed = true
				m.userIDInput.Blur()
			}
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

	if m.userIDInput.Focused() {
		var cmd tea.Cmd
		m.userIDInput, cmd = m.userIDInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View returns the view of the model
func (m *getOrdersModel) View() string {
	if m.changed {
		paginationOpts, err := abstractions.NewPaginationOptions(
			abstractions.WithPage(m.page),
			abstractions.WithPageSize(m.pageSize),
		)
		if err != nil {
			return err.Error()
		}
		orders, err := m.useCase.GetOrders(
			m.userID,
			abstractions.WithPaginationOptions(paginationOpts),
		)
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

	if m.userIDInput.Focused() {
		return m.userIDInput.View()
	}

	return m.table.View() + "\n" + "Page: " + strconv.Itoa(m.page+1) + "\n" + m.table.HelpView()
}
