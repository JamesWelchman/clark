/*
protocol is the specification of the i3bar protocol.
The Block type is what we send to i3bar.
The Click type is what we receive from i3bar.

See: https://i3wm.org/docs/i3bar-protocol.html
*/
package protocol

type (
	Block struct {
		FullText   string `json:"full_text"`
		ShortText  string `json:"short_text,omitempty"`
		Color      string `json:"color,omitempty"`
		Background string `json:"background,omitempty"`
		Border     string `json:"border,omitempty"`
		MinWidth   int    `json:"min_width,omitempty"`
		Align      string `json:"align,omitempty"`
		Name       string `json:"name,omitempty"`
		Instance   string `json:"instance,omitempty"`
		Urgent     bool   `json:"urgent,omitempty"`
		Separator  bool   `json:"separator,omitempty"`
		Markup     string `json:"markup,omitempty"`
	}

	Click struct {
		Name      string   `json:"name"`
		Instance  string   `json:"instance,omitempty"`
		Modifiers []string `json:"modifiers,omitempty"`
		X         int      `json:"x,omitempty"`
		Y         int      `json:"y,omitempty"`
		Button    int      `json:"button,omitempty"`
		RelativeX int      `json:"relative_x,omitempty"`
		RelativeY int      `json:"relative_y,omitempty"`
		Width     int      `json:"width,omitempty"`
		Height    int      `json:"height,omitempty"`
	}
)
