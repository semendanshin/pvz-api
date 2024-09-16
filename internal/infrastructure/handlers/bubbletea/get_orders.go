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

const (
	getOrdersFormUserIDInput = iota
	getOrdersFormLastNInput
	getOrdersFormSamePVZInput
)

func initGetOrdersFormSamePVZInputInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	inputs[getOrdersFormUserIDInput] = textinput.New()
	inputs[getOrdersFormUserIDInput].Focus()
	inputs[getOrdersFormUserIDInput].Prompt = "User ID: "
	inputs[getOrdersFormUserIDInput].Placeholder = "Enter user ID"

	inputs[getOrdersFormLastNInput] = textinput.New()
	inputs[getOrdersFormLastNInput].Prompt = "Last N: "
	inputs[getOrdersFormLastNInput].Placeholder = "Enter last N"

	inputs[getOrdersFormSamePVZInput] = textinput.New()
	inputs[getOrdersFormSamePVZInput].Prompt = "Same PVZ(y/n): "
	inputs[getOrdersFormSamePVZInput].Placeholder = "Enter y/n or leave empty"

	return inputs
}

type validatedInput struct {
	userID  string
	lastN   int
	samePVZ bool
}

func processUserIDInput(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("userID is empty")
	}

	return value, nil
}

func processLastNInput(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	lastN, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("lastN is invalid")
	}

	if lastN < 0 {
		return 0, fmt.Errorf("lastN is negative")
	}

	return lastN, nil
}

func processSamePVZInput(value string) (bool, error) {
	if value == "y" {
		return true, nil
	}

	if value == "n" || value == "" {
		return false, nil
	}

	return false, fmt.Errorf("samePVZ is invalid")
}

func processInput(values []string) (validatedInput, error) {
	userIDValue := values[getOrdersFormUserIDInput]
	lastNValue := values[getOrdersFormLastNInput]
	samePVZValue := values[getOrdersFormSamePVZInput]

	var err error
	var input validatedInput

	input.userID, err = processUserIDInput(userIDValue)
	if err != nil {
		return validatedInput{}, err
	}

	input.lastN, err = processLastNInput(lastNValue)
	if err != nil {
		return validatedInput{}, err
	}

	input.samePVZ, err = processSamePVZInput(samePVZValue)
	if err != nil {
		return validatedInput{}, err
	}

	return input, nil
}

func getOrdersFormSubmitFunc(o *getOrdersModel) func(values []string) error {
	return func(values []string) error {
		input, err := processInput(values)
		if err != nil {
			return err
		}

		o.userID = input.userID
		o.lastN = input.lastN
		o.samePVZ = input.samePVZ

		o.changed = true
		o.settingsFormActive = false

		return nil
	}
}

func initFormModel(o *getOrdersModel) *FormModel {

	inputs := initGetOrdersFormSamePVZInputInputs()

	submit := getOrdersFormSubmitFunc(o)

	return NewFormModel(inputs, submit)
}

// Init initializes the model
func (m *getOrdersModel) Init() tea.Cmd {
	return nil
}

func (m *getOrdersModel) paginateDown() {
	if len(m.data) > 1 {
		m.cursorHistory = append(m.cursorHistory, m.cursor)
		m.cursor = m.data[1].OrderID
		m.changed = true
	}
}

func (m *getOrdersModel) paginateUp() {
	if len(m.data) != 0 {
		if len(m.cursorHistory) != 0 {
			m.cursor = m.cursorHistory[len(m.cursorHistory)-1]
			m.cursorHistory = m.cursorHistory[:len(m.cursorHistory)-1]
			m.changed = true
		}
	}
}

func (m *getOrdersModel) handleKeyboardLetters(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j":
		m.paginateDown()
	case "k":
		m.paginateUp()
	}

	return nil
}

func (m *getOrdersModel) handleKeyboard(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		m.settingsFormActive = true
	case tea.KeyDown:
		m.paginateDown()
	case tea.KeyUp:
		m.paginateUp()
	default:
		m.handleKeyboardLetters(msg)
	}

	return nil
}

func (m *getOrdersModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	switch tea.MouseEvent(msg).Button {
	case tea.MouseButtonWheelDown:
		m.paginateDown()
	case tea.MouseButtonWheelUp:
		m.paginateUp()
	default:
	}

	return nil
}

func (m *getOrdersModel) innerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			strconv.Itoa(order.Weight),
			strconv.Itoa(order.Cost),
			order.Packaging.String(),
			strconv.FormatBool(order.AdditionalFilm),
		}
	}
	m.table.SetRows(rows)
	m.data = orders

	return nil
}

func (m *getOrdersModel) innerView() string {
	if m.changed {
		if err := m.updateData(); err != nil {
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
