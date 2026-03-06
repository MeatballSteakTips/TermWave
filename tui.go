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
	savedStations   []Station //Page number, each will hold 8
	stationCursor   int
	viewState 		  string // Tells me if I am in search or saved sations mode
	savedPage       int    //Page #s
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

	//Loading saved stations
	loadedStations, err := loadStations()
	if err != nil || loadedStations == nil {
		loadedStations = []Station{}
	}

	return model {
		stations: []Station{},
		savedStations: loadedStations,
		stationCursor: 0,
		focused: "stations",
		viewState: "saved",
		savedPage: 0,
		searchInput: ti,
		activeMenuIndex: 0,
		menuOpen: false,
		menuCursors: []int{0, 0, 0}, 
		menuTitles: []string{"Stations", "Settings", "Help"},
		menuItems: [][]string {
						{"Add Station","Saved Stations", "Quit"},
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
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stationsLoadedMsg:
		//Making a loop so I can check if the station is already saved
		for i, newStation := range msg {
			for _, savedStation := range m.savedStations {
				if newStation.URL == savedStation.URL {
					msg[i].Saved = "*"
					break
				}
			}
		}
		m.stations = msg
		m.viewState = "search"
		m.stationCursor = 0
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyPressMsg:
		s := msg.String();
		if m.focused == "search" { //This is for the search window. Note to me: I should make this less confusing later
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
			//This is to set up the field for typing. Otherwise it will do weird stuff
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
				} else if selectedItem == "Saved Stations" {
					m.focused = "stations"
					m.viewState = "saved"
					m.menuOpen = false
				} else if selectedItem == "Quit" {
					StopStream()
					return m, tea.Quit
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
				var currentListLen int
				if m.viewState == "search" {
					currentListLen = len(m.stations)
				} else {
					startIndex := m.savedPage * 16
					endIndex := startIndex + 16
					if endIndex > len(m.savedStations) {
						endIndex = len(m.savedStations)
					}
					currentListLen = endIndex - startIndex
				}
				switch s {
				case "up", "k":
					if m.stationCursor > 0 {
						m.stationCursor--
					}
				case "down", "j":
					if m.stationCursor < currentListLen - 1 {
						m.stationCursor++
					}
				case "s":
					if m.viewState == "search" && len(m.stations) > 0 && m.stations[m.stationCursor].Saved != "*" {
						m.savedStations = append(m.savedStations, m.stations[m.stationCursor])
						m.stations[m.stationCursor].Saved = "*"
						_ = saveStations(m.savedStations)
					}
				case "x", "delete":
					if m.viewState == "saved" && len(m.savedStations) > 0 {
						startIndex := m.savedPage * 16
						actualIndex := startIndex + m.stationCursor

						m.savedStations = append(m.savedStations[:actualIndex], m.savedStations[actualIndex + 1:]...)
						_ = saveStations(m.savedStations)

						itemsLeftOnPage := len(m.savedStations) - startIndex
							if m.stationCursor >0 && m.stationCursor >= itemsLeftOnPage {
								m.stationCursor--
							}

							if startIndex >= len(m.savedStations) && m.savedPage > 0 {
								m.savedPage--
								m.stationCursor = 7
							}
					}
				case "enter":
					if m.viewState == "search" && len(m.stations) > 0 {
						m.currentStation = m.stations[m.stationCursor]
						_ = PlayStream(m.currentStation.URL)
					} else if m.viewState == "saved" && currentListLen > 0 {
						startIndex := m.savedPage * 16
						m.currentStation = m.savedStations[startIndex + m.stationCursor]
						_ = PlayStream(m.currentStation.URL)
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

	if m.viewState == "search" {
		leftContent = "Stations (Search Results)\n\n"
		if m.err != nil {
			leftContent += fmt.Sprintf("Error fetching stations: %v", m.err)
		} else if len(m.stations) == 0 {
			leftContent += "Loading stations..."
		} else {
			for i, s := range m.stations {
				cursor := "  "
				if m.stationCursor == i && m.focused == "stations" {
					cursor = "> "
				}
				leftContent += fmt.Sprintf("%s%d. %s %s\n", cursor, i + 1, s.Name, s.Saved)
			}
		}
	} else if m.viewState == "saved" { //Saved Stations logic
		leftContent = "Stations\n\n"

		if len(m.savedStations) == 0 {
			leftContent += "No stations saved\nSearch for a station in Stations->Add Station"
		} else {
			itemsPerPage := 16

			startIndex := m.savedPage * itemsPerPage
			endIndex := startIndex + itemsPerPage

			if endIndex > len(m.savedStations) {
				endIndex = len(m.savedStations)
			}

			pageItems := m.savedStations[startIndex:endIndex]

			for i, s := range pageItems {
				cursor := "  "

				if m.stationCursor == i && m.focused == "stations" {
					cursor = "> "
				}
				leftContent += fmt.Sprintf("%s%d. %s %s\n", cursor, (startIndex + i) + 1, s.Name, s.Saved)
			}

			totalPages := (len(m.savedStations) + itemsPerPage - 1) / itemsPerPage
			leftContent += fmt.Sprintf("\n\n  --- Page %d of %d ---", m.savedPage + 1, totalPages)
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
