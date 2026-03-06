# TermWave
### A Terminal-based internet radio player for Linux

Termwave is a TUI that searches for internet radio stations and plays them back. It is my first program made in Go, and the first using Bubbletea, so it is still in early-stage development. 

⚠️ Termwave is currently still being developed. While I will be actively adding/fixing features, things may break as I add/fix features.

Program was made using Go and the BubbleTea V2 framework. ⚠️ MPV is required for audio playback!

### Currently Working:
    * Search for internet radio stations (Using api.radio-browser)
    * Radio Playback using mpv
    * Functionality to Save and Remove Stations has been added
### Roadmap:
    * Adding live song-titles/ICY metadata
    * Adding playback-buttons
    * Make the UI look better
Will be making releases hopefully in the near-future

### Usage:

```
$ git clone https://github.com/MeatballSteakTips/TermWave.git
$ cd TermWave
$ go run .
```
### Controls: 
    * Use Tab to switch focus between Toolbar and Stations Pane.
    * S saves a Station
    * X removes a station
    * Enter for everything else

Be sure you have MPV installed!

Debian/Ubuntu based:
```
$ sudo apt install mpv
```

Arch
```
sudo pacman -S mpv
```

Thank you for trying it out! I hope to make it better in the near future!


