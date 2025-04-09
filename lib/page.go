package lib

type Page string

const (
	PAGE_HOME         Page = "home"
	PAGE_SETUP_WIZARD Page = "setup-wizard"
)

type PanelPage string

const (
	PANEL_PAGE_SONGS     PanelPage = "songs"
	PANEL_PAGE_PLAYLISTS PanelPage = "playlists"
	PANEL_PAGE_SETTINGS  PanelPage = "settings"
)
