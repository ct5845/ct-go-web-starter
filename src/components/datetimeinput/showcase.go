package datetimeinput

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	_ "embed"
	"html/template"
)

var (
	//go:embed showcase.js
	showcaseJS string
	// showcaseScriptTpl carries only the field-change-linking script used by
	// renderDateRange's demo — nothing else in this showcase needs it, so it
	// stays local to this file rather than living in the datetimeinput
	// package proper.
	showcaseScriptTpl = component.WithAlpine("showcase.html", "", showcaseJS)
)

var Showcase = demo.Page{
	Slug:        "datetimeinput",
	Title:       "Date/time input",
	Source:      "components/datetimeinput",
	Description: "Date, time, and combined date+time inputs, using the browser's own picker.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	date, err := Render(Options{
		Id:       "showcase-datetimeinput-date",
		Name:     "birthday",
		Label:    "Birthday",
		Type:     Date,
		Max:      "2026-08-18",
		Hint:     "Must be in the past.",
		Required: true,
	})
	if err != nil {
		return "", err
	}

	time, err := Render(Options{
		Id:    "showcase-datetimeinput-time",
		Name:  "alarm",
		Label: "Alarm",
		Type:  Time,
		Step:  "300",
		Hint:  "In 5-minute steps.",
	})
	if err != nil {
		return "", err
	}

	dateTime, err := Render(Options{
		Id:    "showcase-datetimeinput-datetime",
		Name:  "appointment",
		Label: "Appointment",
		Type:  DateTimeLocal,
		Hint:  "Local time.",
	})
	if err != nil {
		return "", err
	}

	dateRange, err := renderDateRange()
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Each renders the browser's own date/time picker. Birthday sets max to today and is required; Alarm steps in 5-minute increments.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: date + time + dateTime + dateRange,
	})
}

// renderDateRange demonstrates the field-change event every field's wrapper
// dispatches automatically via inputControl (see field.go's doc comment):
// dragging Start fires field-change with its id and new value, and
// showcaseDateRange (registered in showcase.js) shifts End by the same
// number of days, preserving the gap between the two dates. Neither field
// knows the other exists — the wiring lives entirely in this showcase, not
// in datetimeinput itself.
const (
	showcaseRangeStartId = "showcase-datetimeinput-range-start"
	showcaseRangeEndId   = "showcase-datetimeinput-range-end"
)

func renderDateRange() (template.HTML, error) {
	script, err := showcaseScriptTpl.Render(nil)
	if err != nil {
		return "", err
	}

	start, err := Render(Options{
		Id:    showcaseRangeStartId,
		Name:  "trip-start",
		Label: "Trip start",
		Type:  Date,
		Value: "2026-09-01",
		Hint:  "Dragging this keeps the trip length fixed.",
	})
	if err != nil {
		return "", err
	}

	end, err := Render(Options{
		Id:    showcaseRangeEndId,
		Name:  "trip-end",
		Label: "Trip end",
		Type:  Date,
		Value: "2026-09-08",
	})
	if err != nil {
		return "", err
	}

	wrapper := `<div x-data="showcaseDateRange('` + showcaseRangeStartId + `', '` + showcaseRangeEndId + `')" x-on:field-change.window="shiftEnd($event)" class="flex flex-col gap-4">` +
		string(start) + string(end) + `</div>`

	return demo.Frame(demo.FrameOptions{
		Note:    "A cross-field example: Start and End are otherwise independent native inputs. showcaseDateRange listens for Start's field-change event and shifts End by the same delta, so the trip length survives a drag on Start.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: script + template.HTML(wrapper),
	})
}
