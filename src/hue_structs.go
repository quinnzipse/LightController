package main

type Response struct {
	Errors []string     `json:"errors"`
	Data   []LightState `json:data`
}

type LightState struct {
	Identifier   *string    `json:"id,omitempty"`
	IdentifierV1 *string    `json:"id_v1,omitempty"`
	On           *On        `json:"on,omitempty"`
	Metadata     *Metadata  `json:"metadata,omitempty"`
	Color        *Color     `json:"color,omitempty" default:nil`
	ColorTemp    *ColorTemp `json:"color_temperature,omitempty"`
	Dimming      *Dimming   `json:"dimming,omitempty"`
}

type On struct {
	Value bool `json:"on"`
}

type Color struct {
	Gamut     *Gamut       `json:"gamut,omitempty"`
	GamutType string       `json:"gamut_type,omitempty"`
	XY        *Coordinates `json:"xy,omitempty"`
}

type Gamut struct {
	Blue  *Coordinates `json:"blue,omitempty"`
	Green *Coordinates `json:"green,omitempty"`
	Red   *Coordinates `json:"red,omitempty"`
}

type Coordinates struct {
	X float32 `json:"x,omitempty"`
	Y float32 `json:"y,omitempty"`
}

type ColorTemp struct {
	Mirek       int          `json:"mirek,omitempty"`
	MirekSchema *MirekSchema `json:"mirek_schema,omitempty"`
	MirekValid  bool         `json:"mirek_valid,omitempty"`
}

type MirekSchema struct {
	Maximum int `json:"mirek_maximum,omitempty`
	Minimum int `json:"mirek_minimum,omitempty`
}

type Dimming struct {
	Brightness  uint8 `json:"brightness,omitempty"`
	MinDimLevel uint8 `json:"min_dim_level,omitempty"`
}

type Metadata struct {
	Type string `json:"archetype,omitempty"`
	Name string `json:"name,omitempty"`
}
