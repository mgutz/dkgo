package dk

type Status struct {
	Clients    []*Client    `json:"clients"`
	Workspaces []*Workspace `json:"workspaces"`
	Monitors   []*Monitor   `json:"monitors"`
	Rules      []*Rule      `json:"rules"`
	Global     *Global      `json:"global"`
	Panels     []*Panel     `json:"panels"`
}

type Panel struct {
	ID       string   `json:"id"`
	Class    string   `json:"class"`
	Instance string   `json:"instance"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	W        int      `json:"w"`
	H        int      `json:"h"`
	L        int      `json:"l"`
	R        int      `json:"r"`
	T        int      `json:"t"`
	B        int      `json:"b"`
	Monitor  *Monitor `json:"monitor"`
}

type Workspace struct {
	Name       string    `json:"name"`
	Number     int       `json:"number"`
	Focused    bool      `json:"focused"`
	Monitor    string    `json:"monitor"`
	Layout     string    `json:"layout"`
	Master     int       `json:"master"`
	Stack      int       `json:"stack"`
	MSplit     float64   `json:"msplit"`
	SSplit     float64   `json:"ssplit"`
	Gap        int       `json:"gap"`
	SmartGap   bool      `json:"smart_gap"`
	PadL       int       `json:"pad_l"`
	PadR       int       `json:"pad_r"`
	PadT       int       `json:"pad_t"`
	PadB       int       `json:"pad_b"`
	Clients    []*Client `json:"clients"`
	FocusStack []*Client `json:"focus_stack"`
	// used for dynamic workspaces
	_skip bool
}

type Global struct {
	NumWS       int      `json:"numws"`
	StaticWS    bool     `json:"static_ws"`
	FocusMouse  bool     `json:"focus_mouse"`
	FocusOpen   bool     `json:"focus_open"`
	FocusUrgent bool     `json:"focus_urgent"`
	WinMinWH    int      `json:"win_minwh"`
	WinMinXY    int      `json:"win_minxy"`
	SmartBorder bool     `json:"smart_border"`
	SmartGap    bool     `json:"smart_gap"`
	TileHints   bool     `json:"tile_hints"`
	TileToHead  bool     `json:"tile_tohead"`
	ObeyMotif   bool     `json:"obey_motif"`
	Layouts     []string `json:"layouts"`
	Callbacks   []string `json:"callbacks"`
	Border      struct {
		Width        int    `json:"width"`
		OuterWidth   int    `json:"outer_width"`
		Focus        string `json:"focus"`
		Urgent       string `json:"urgent"`
		Unfocus      string `json:"unfocus"`
		OuterFocus   string `json:"outer_focus"`
		OuterUrgent  string `json:"outer_urgent"`
		OuterUnfocus string `json:"outer_unfocus"`
	} `json:"border"`
	Focused *Monitor `json:"focused"`
}

type Monitor struct {
	Name      string     `json:"name"`
	Number    int        `json:"number"`
	Focused   bool       `json:"focused"`
	X         int        `json:"x"`
	Y         int        `json:"y"`
	W         int        `json:"w"`
	H         int        `json:"h"`
	WX        int        `json:"wx"`
	WY        int        `json:"wy"`
	WW        int        `json:"ww"`
	WH        int        `json:"wh"`
	Workspace *Workspace `json:"workspace"`
}

type Rule struct {
	Title     string `json:"title"`
	Class     string `json:"class"`
	Instance  string `json:"instance"`
	Workspace int    `json:"workspace"`
	Monitor   string `json:"monitor"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	W         int    `json:"w"`
	H         int    `json:"h"`
	Float     bool   `json:"float"`
	Full      bool   `json:"full"`
	FakeFull  bool   `json:"fakefull"`
	Sticky    bool   `json:"sticky"`
	Focus     bool   `json:"focus"`
	IgnoreCfg bool   `json:"ignore_cfg"`
	IgnoreMsg bool   `json:"ignore_msg"`
	Callback  string `json:"callback"`
	XGrav     string `json:"xgrav"`
	YGrav     string `json:"ygrav"`
}

type Client struct {
	ID        string  `json:"id"`
	PID       int     `json:"pid"`
	Title     string  `json:"title"`
	Class     string  `json:"class"`
	Instance  string  `json:"instance"`
	Workspace int     `json:"workspace"`
	Focused   bool    `json:"focused"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	W         int     `json:"w"`
	H         int     `json:"h"`
	BW        int     `json:"bw"`
	Hoff      int     `json:"hoff"`
	Float     bool    `json:"float"`
	Full      bool    `json:"full"`
	FakeFull  bool    `json:"fakefull"`
	Sticky    bool    `json:"sticky"`
	Urgent    bool    `json:"urgent"`
	Above     bool    `json:"above"`
	Hidden    bool    `json:"hidden"`
	Scratch   bool    `json:"scratch"`
	Callback  string  `json:"callback"`
	TransID   string  `json:"trans_id"`
	Absorbed  *Client `json:"absorbed"`
}
