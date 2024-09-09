package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"reflect"
	"strconv"
)

var _ tea.Model = &getOrdersModel{}

// getOrdersModel is a model for getting orders
type getOrdersModel struct {
	useCase abstractions.IPVZOrderUseCase

	settingsForm       *FormModel
	settingsFormActive bool

	table table.Model

	userID  string
	lastN   int
	samePVZ bool

	data    []domain.PVZOrder
	changed bool

	cursor        string
	cursorHistory []string
	pageSize      int
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
		{Title: "Weight", Width: 10},
		{Title: "Cost", Width: 10},
		{Title: "Packaging", Width: 10},
		{Title: "AdditionalFilm", Width: 15},
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

	model := &getOrdersModel{
		useCase: useCase,

		settingsFormActive: true,
		table:              dataTable,

		data: make([]domain.PVZOrder, 0),

		cursorHistory: make([]string, 0),
		pageSize:      pageSize,
	}

	model.settingsForm = initFormModel(model)

	return model
}

func initFormModel(o *getOrdersModel) *FormModel {
	const (
		userIDInput = iota
		lastNInput
		samePVZInput
	)

	inputs := make([]textinput.Model, 3)

	inputs[userIDInput] = textinput.New()
	inputs[userIDInput].Focus()
	inputs[userIDInput].Prompt = "User ID: "
	inputs[userIDInput].Placeholder = "Enter user ID"

	inputs[lastNInput] = textinput.New()
	inputs[lastNInput].Prompt = "Last N: "
	inputs[lastNInput].Placeholder = "Enter last N"

	inputs[samePVZInput] = textinput.New()
	inputs[samePVZInput].Prompt = "Same PVZ(y/n): "
	inputs[samePVZInput].Placeholder = "Enter y/n or leave empty"

	submit := func(values []string) error {
		userIDValue := values[userIDInput]
		lastNValue := values[lastNInput]
		samePVZValue := values[samePVZInput]

		var err error

		var input struct {
			userID  string
			lastN   int
			samePVZ bool
		}
		{
			if userIDValue == "" {
				return fmt.Errorf("userID is empty")
			}

			if lastNValue != "" {
				input.lastN, err = strconv.Atoi(lastNValue)
				if err != nil {
					return fmt.Errorf("lastN is invalid")
				}
			}

			if input.lastN < 0 {
				return fmt.Errorf("lastN is negative")
			}

			if samePVZValue != "y" && samePVZValue != "n" && samePVZValue != "" {
				return fmt.Errorf("samePVZ is invalid")
			}

			if samePVZValue == "y" {
				input.samePVZ = true
			} else {
				input.samePVZ = false
			}

			input.userID = userIDValue
		}

		o.userID = input.userID
		o.lastN = input.lastN
		o.samePVZ = input.samePVZ

		o.changed = true
		o.settingsFormActive = false

		return nil
	}

	return NewFormModel(inputs, submit)
}

// Init initializes the model
func (m *getOrdersModel) Init() tea.Cmd {
	return nil
}

func (m *getOrdersModel) innerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	paginateDown := func() {
		if len(m.data) > 1 {
			m.cursorHistory = append(m.cursorHistory, m.cursor)
			m.cursor = m.data[1].OrderID
			m.changed = true
		}
	}

	paginateUp := func() {
		if len(m.data) != 0 {
			if len(m.cursorHistory) != 0 {
				m.cursor = m.cursorHistory[len(m.cursorHistory)-1]
				m.cursorHistory = m.cursorHistory[:len(m.cursorHistory)-1]
				m.changed = true
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.settingsFormActive = true
		case tea.KeyDown:
			paginateDown()
		case tea.KeyUp:
			paginateUp()
		default:
			switch msg.String() {
			case "j":
				paginateDown()
			case "k":
				paginateUp()
			}
		}
	case tea.MouseMsg:
		switch tea.MouseEvent(msg).Button {
		case tea.MouseButtonWheelDown:
			paginateDown()
		case tea.MouseButtonWheelUp:
			paginateUp()
		default:
		}
	default:
	}

	return m, nil
}

// Update updates the model
func (m *getOrdersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.settingsFormActive {
		return m.innerUpdate(msg)
	}

	_, cmd := m.settingsForm.Update(msg)
	if cmd != nil {
		if reflect.ValueOf(cmd).Pointer() == reflect.ValueOf(tea.Quit).Pointer() && !m.settingsFormActive {
			m.settingsFormActive = false
			cmd = nil
		}
	}

	return m, cmd
}

func (m *getOrdersModel) updateData() error {
	opts := []abstractions.GetOrdersOptFunc{
		abstractions.WithCursorID(m.cursor),
		abstractions.WithLimit(m.pageSize),
		abstractions.WithLastNOrders(m.lastN),
	}
	if m.samePVZ {
		opts = append(opts, abstractions.WithSamePVZ())
	}
	orders, err := m.useCase.GetOrders(
		m.userID,
		opts...,
	)
	if err != nil {
		return err
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
	m.data = orders

	return nil
}

func (m *getOrdersModel) innerView() string {
	if m.changed {
		err := m.updateData()
		if err != nil {
			return err.Error()
		}
		m.changed = false
	}

	return m.table.View() + "\n\n" + m.table.HelpView()
}

// View returns the view
func (m *getOrdersModel) View() string {
	if m.settingsFormActive {
		return m.settingsForm.View()
	}

	return m.innerView()
}
