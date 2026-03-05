package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"charm.land/bubbles/v2/textinput"
)

type stationsLoadedMsg []Station
type errMsg struct{ err error }

type model struct {
	width       		int
	height 					int
	activeMenuIndex int
	menuOpen   	 		bool
	menuCursors 		[]int
	menuTitles  		[]string
	menuItems   		[][]string
	stations 				[]Station
	stationCursor   int
	focused         string
	err							error
	currentStation  Station
	searchInput     textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search Station by Name..."
	ti.CharLimit = 50
	ti.SetWidth(30)

	return model {
		stations: []Station{},
		stationCursor: 0,
		focused: "stations",
		searchInput: ti,
		activeMenuIndex: 0,
		menuOpen: false,
		menuCursors: []int{0, 0, 0}, 
		menuTitles: []string{"Stations", "Settings", "Help"},
		menuItems: [][]string {
						{"Add Station", "Remove Station", "Quit"},
						{"Audio Settings", "Theme Settings", "Preferences"},
						{"About", "Documentation", "License"},
		},
	}
}

func fetchStations(query string) tea.Cmd {
	return func() tea.Msg {
		stations, err := StationSearch(query)
		if err != nil {
			return errMsg{err}
		}
		return stationsLoadedMsg(stations)
	}
}

func (m model) Init() tea.Cmd {
	return fetchStations("Synphaera")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stationsLoadedMsg:
		m.stations = msg
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyPressMsg:
		s := msg.String();
		if m.focused == "search" {
			switch s {
			case "esc":
				m.focused = "stations"
				m.searchInput.Blur()
				return m, nil
			case "enter":
				query := m.searchInput.Value()
				m.focused = "stations"
				m.searchInput.Blur()
				return m, fetchStations(query) //Search for the list

			}
			//This is to set up the field for typing
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd

		}
		if s == "ctrl+c" || s == "q" || s == "esc" {
			StopStream()
			return m, tea.Quit
		}

		if s == "tab" && !m.menuOpen {
			if m.focused == "stations" {
				m.focused = "toolbar"
			} else {
				m.focused = "stations"
			}
			return m, nil
		}
		if m.menuOpen {
			switch s {
			case "up", "k":
				if m.menuCursors[m.activeMenuIndex] > 0 {
					m.menuCursors[m.activeMenuIndex]--
				}
			case "down", "j":
				maxIndex := len(m.menuItems[m.activeMenuIndex]) - 1
				if m.menuCursors[m.activeMenuIndex] < maxIndex {
					m.menuCursors[m.activeMenuIndex]++
				}
			case "esc", "q", "left", "right":
				m.menuOpen = false
			case "enter":
				selectedItem := m.menuItems[m.activeMenuIndex][m.menuCursors[m.activeMenuIndex]]

				if selectedItem == "Add Station" {
					m.focused = "search"
					m.searchInput.Focus()
					m.searchInput.SetValue("")
					m.menuOpen = false
				} else {
					m.menuOpen = false
				}
			}

		} else {
			// Menu is closed
			if m.focused == "toolbar" {
				switch s {
				case "right", "l":
					m.activeMenuIndex = (m.activeMenuIndex + 1) % len(m.menuTitles)
				case "left", "h":
					m.activeMenuIndex = ((m.activeMenuIndex - 1) + len(m.menuTitles)) % len(m.menuTitles)
				case "enter", "down", "s":
					m.menuOpen = true
				case "q", "esc":
					StopStream()
					return m, tea.Quit
				}
				
			} else if m.focused == "stations" {
				//Radio list focused
				switch s {
				case "up", "k":
					if m.stationCursor > 0 {
						m.stationCursor--
					}
				case "down", "j":
					if m.stationCursor < len(m.stations)-1 {
						m.stationCursor++
					}
				case "enter":
					//set current station and play
					m.currentStation = m.stations[m.stationCursor]
					if len(m.stations) > 0 {
						_ = PlayStream(m.stations[m.stationCursor].URL)
					}
				case "q", "esc":
					StopStream()
					return m, tea.Quit
				}
			}
		}

		return m, tea.RequestWindowSize

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}	
	
func (m model) drawPanes() string {
	availW := m.width - 6
	leftW := int(float64(availW) * 0.6)
	rightW := availW - leftW
	paneH := m.height - 5
	leftTitle := "Radios"
	leftContent := fmt.Sprintf("%s\n\n", leftTitle)
	stationName := "None"
	//stationImage := ""

	if m.err != nil {
		leftContent += fmt.Sprintf("Error fetching stations: %v", m.err)
	} else if len(m.stations) == 0 {
		leftContent += "Loading stations..."
	} else {
		for i, s := range m.stations {
			cursor := " "
			if m.stationCursor == i {
				if m.focused == "stations" {
					cursor = "> "
				} else {
					cursor = "  "
				}
			}
			leftContent += fmt.Sprintf("%s%d. %s\n", cursor, i + 1, s.Name)
		}
	}

	if m.currentStation.Name != "" {
		stationName = m.currentStation.Name
		//stationImage = m.currentStation.Image
	}

	rightContent := fmt.Sprintf("Now Playing\n\nStation: %s\nTitle: ", stationName)
	leftPane  := paneStyle.Width(leftW).Height(paneH).Render(leftContent)
	rightPane := paneStyle.Width(rightW).Height(paneH).Render(rightContent)
	
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func (m model) drawToolbar(height int) string {
	var buttons []string

	for i, title := range m.menuTitles {
		activeStyle := buttonStyle

		if i == m.activeMenuIndex && m.focused == "toolbar" {
			if m.menuOpen {
				activeStyle = activeStyle.Background(lipgloss.Color("62"))
			} else {
				activeStyle = activeStyle.Background(lipgloss.Color("240"))
			}
		} else {
			activeStyle = activeStyle.Background(lipgloss.Color("0"))
		}
		buttons = append(buttons, activeStyle.Render(fmt.Sprintf(" %s ", title)))
	}
	toolbarContent := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	toolbar := toolbarStyle.Width(m.width - 2).
					Width(m.width - 2).
					Height(height).
					Render(toolbarContent)
	
	return toolbar
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing")
	}
	toolbarHeight := 1
	contentPanes := m.drawPanes()  
	toolbar := m.drawToolbar(toolbarHeight)

  border := baseBorderStyle.
					Width(m.width - 2).
					Height(m.height - toolbarHeight - 4).
					Render(contentPanes)

	ui := lipgloss.JoinVertical(lipgloss.Left, toolbar, border)
	bgLayer := lipgloss.NewLayer(ui).X(0).Y(0).Z(0)
	layers := []*lipgloss.Layer{bgLayer}

	if m.menuOpen {
		currentItems := m.menuItems[m.activeMenuIndex]
		currentCursor := m.menuCursors[m.activeMenuIndex]

		menuText := fmt.Sprintf("%s\n\n", m.menuTitles[m.activeMenuIndex])
		for i, item := range currentItems {
			cursor := " "
			if currentCursor == i {
				cursor = ">"
			}
			menuText += fmt.Sprintf(" %s %s\n", cursor, item)
		}

		menu := menuStyle.Background(lipgloss.Color("0")).Render(menuText)
	
		xOffset := 2 + m.activeMenuIndex * 14
		menuLayer := lipgloss.NewLayer(menu).X(xOffset).Y(3).Z(1)

		layers = append(layers, menuLayer)
	}

	if m.focused == "search" {
		popupContent := fmt.Sprintf("Search Station by Name:\n\n%s", m.searchInput.View())

		searchPopup := popup.Render(popupContent)
		popupWidth := lipgloss.Width(searchPopup)
		popupHeight := lipgloss.Height(searchPopup)
		x := (m.width / 2) - (popupWidth / 2)
		y := (m.height / 2) - (popupHeight / 2)

		searchLayer := lipgloss.NewLayer(searchPopup).X(x).Y(y).Z(2)
		layers = append(layers, searchLayer)
	}
	compositor := lipgloss.NewCompositor(layers...)

	v := tea.NewView(compositor.Render())
  v.AltScreen = true
	return v
}
