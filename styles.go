package main

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
		baseBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Center)

		toolbarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Top)

		buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("0")).
			Padding(0, 1).
			Bold(true)

		menuStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			BorderBackground(lipgloss.Color("0")).
			Padding(0, 1)
			//.Align(lipgloss.Center)
		
		paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)		
			
		popup = baseBorderStyle.
			BorderForeground(lipgloss.Color("62")). // Make it pop with a purple border
			Padding(1, 2)
)

