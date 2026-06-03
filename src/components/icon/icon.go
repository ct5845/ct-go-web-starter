package icon

import (
	"fmt"
	"sort"
	"strings"
)

const (
	CalendarMonth = "calendar_month"
	Home          = "home"
	Close         = "close"
	Menu          = "menu"
	Search        = "search"
	Colors        = "colors"
	Serif         = "serif"
	RightArrow    = "keyboard_arrow_right"
)

var IconFontHref string

func init() {
	var all = []string{
		CalendarMonth,
		Home,
		Close,
		Colors,
		Serif,
		Menu,
		Search,
		RightArrow,
	}
	sort.Strings(all)

	IconFontHref = fmt.Sprintf("https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&icon_names=%s&display=block", strings.Join(all, ","))
}
