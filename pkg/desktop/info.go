package desktop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mgutz/dkgo/pkg/util"
	"gopkg.in/ini.v1"
)

var (
	desktopFilesCache = make(map[string][]*AppInfo)
	iconCache         = make(map[string]string)
)

// AppInfo represents the information of a desktop file.
type AppInfo struct {
	Name           string
	Exec           string
	StartupWMClass string
	Icon           string
	Terminal       bool
	Categories     []string
	Keywords       []string
	Bin            string
}

// GetIcon guesses the icon by reading desktop files. It uses the WM_CLASS
// property of the window to find the icon. If the icon is not found, it then tries
// the commandline to find the icon.
func GetIcon(wmClass string, commandline string) (string, error) {
	// check if the icon is already in the cache
	if icon, ok := iconCache[wmClass]; ok {
		return icon, nil
	}

	clientBin := ""

	if commandline != "" {
		clientBin = filepath.Base(strings.Split(commandline, " ")[0])
	}

	apps, err := ReadDesktopFiles("")
	if err != nil {
		panic(err)
	}

	var foundApp *AppInfo
	for _, app := range apps {
		if wmClass != "" {
			if app.StartupWMClass == wmClass {
				foundApp = app
				break
			}
		}
		if commandline != "" {
			if app.Bin == clientBin {
				foundApp = app
				break
			}
		}
	}

	if foundApp != nil {
		iconCache[wmClass] = foundApp.Icon
		return foundApp.Icon, nil
	}

	return "", nil
}

// GetIconOr returns the icon or the fallback icon if the icon is not found.
func GetIconOr(wmClass string, commandline string, fallback string) string {
	icon, err := GetIcon(wmClass, commandline)
	if err != nil || icon == "" {
		return fallback
	}
	return icon
}

// ReadDesktopFiles reads all desktop files from the $XDG_DATA_DIRS/applications.
func ReadDesktopFiles(xdgDataDirs string) ([]*AppInfo, error) {
	if len(xdgDataDirs) == 0 {
		xdgDataDirs = os.Getenv("XDG_DATA_DIRS")
		if xdgDataDirs == "" {
			xdgDataDirs = "~/.local/share:/usr/share"
		} else {
			xdgDataDirs = "~/.local/share:" + xdgDataDirs
		}
	}

	cacheKey := xdgDataDirs

	if apps, ok := desktopFilesCache[cacheKey]; ok {
		return apps, nil
	}

	dirs := strings.Split(xdgDataDirs, ":")

	// reverse the order so that the user's local applications take
	// precedence
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}

	appInfoList := map[string]*AppInfo{}

	// Read desktop files from the specified directories
	for _, dir := range dirs {
		dir = filepath.Join(util.Untildify(dir), "applications")

		// Check if the directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// get list of *.desktop files from directory
		files, err := filepath.Glob(filepath.Join(dir, "*.desktop"))
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			// parse the file and create an AppInfo object
			info, err := parseDesktopFile(file)
			if err != nil {
				return nil, err
			}

			name := util.BaseNoExt(file)

			if old, ok := appInfoList[name]; ok {
				// merge the two AppInfo objects
				old.Name = info.Name
				old.Exec = info.Exec
				if info.StartupWMClass != "" {
					old.StartupWMClass = info.StartupWMClass
				}
				if info.Icon != "" {
					old.Icon = info.Icon
				}
				if len(info.Categories) > 0 {
					old.Categories = info.Categories
				}
				if len(info.Keywords) > 0 {
					old.Keywords = info.Keywords
				}
				continue
			}

			appInfoList[name] = info
		}
	}

	result := make([]*AppInfo, 0, len(appInfoList))
	for _, v := range appInfoList {
		result = append(result, v)
	}

	desktopFilesCache[cacheKey] = result
	return result, nil
}

func parseDesktopFile(filename string) (*AppInfo, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Parse the desktop file and create an AppInfo object
	section := cfg.Section("Desktop Entry")
	exec := section.Key("Exec").String()
	bin := ""
	if exec != "" {
		bin = strings.Split(exec, " ")[0]
		bin = util.BaseNoExt(bin)
	}

	return &AppInfo{
		Name:           section.Key("Name").String(),
		Exec:           exec,
		StartupWMClass: section.Key("StartupWMClass").String(),
		Icon:           section.Key("Icon").String(),
		Terminal:       section.Key("Terminal").MustBool(),
		Categories:     section.Key("Categories").Strings(";"),
		Keywords:       section.Key("Keywords").Strings(";"),
		Bin:            bin,
	}, nil
}
