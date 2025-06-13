package main

import (
	"fmt"

	tmux "github.com/jubnzv/go-tmux"
)

const DIRPATH = "/home/felix/Projects"
const PROJECT_NAME = "fenix"

type Server struct {
	tmuxServer *tmux.Server
}

func main() {
	exists, config := LoadConfig(CONFIG_FILE)

	if !exists {
		loadSession(PROJECT_NAME, nil)
	} else {
		loadSession(PROJECT_NAME, &config)
	}
}

func loadSession(name string, config *Config) error {
	session, err := createSession(name)
	if err != nil {
		fmt.Println("Error attaching session:", err)
		return err
	}

	if config != nil {
		applyConfig(session, *config)
	} else {
		windows, err := session.ListWindows()
		if err != nil {
			fmt.Println("Error listing windows:", err)
			return err
		}

		for _, window := range windows {
			switchDirectory("", window)
		}
	}

	err = session.AttachSession()
	if err != nil {
		fmt.Println("Error attaching session:", err)
		return err
	}

	return nil
}

func createSession(name string) (tmux.Session, error) {
	server := new(tmux.Server)

	exists, session := getSessionByName(name)

	if exists {
		return session, nil
	}

	session, err := server.NewSession(name)
	if err != nil {
		return tmux.Session{}, fmt.Errorf("error creating session: %w", err)
	}

	return session, nil
}

func applyConfig(session tmux.Session, config Config) error {
	for _, windowConfig := range config.Windows {
		window, err := session.NewWindow(windowConfig.Name)
		if err != nil {
			fmt.Println("Error creating window:", err)
			return err
		}

		switchDirectory(windowConfig.Path, window)

		panes, err := window.ListPanes()
		if err != nil {
			fmt.Println("Error listing panes:", err)
			return err
		}

		err = panes[0].RunCommand(windowConfig.Cmd)
		if err != nil {
			fmt.Println("Error running command:", err)
			return err
		}
	}

	windows, err := session.ListWindows()
	if err != nil {
		fmt.Println("Error listing windows:", err)
		return err
	}

	panes, err := windows[0].ListPanes()
	if err != nil {
		fmt.Println("Error listing panes:", err)
		return err
	}

	err = panes[0].RunCommand("exit")
	if err != nil {
		fmt.Println("Error running command:", err)
		return err
	}

	return nil
}

func switchDirectory(name string, pane tmux.Window) error {
	panes, err := pane.ListPanes()
	if err != nil {
		fmt.Println("Error listing panes:", err)
		return err
	}

	for _, pane := range panes {
		err = pane.RunCommand(fmt.Sprintf("cd %s/%s/%s", DIRPATH, PROJECT_NAME, name))
		if err != nil {
			fmt.Println("Error switching directory:", err)
			return err
		}
	}

	return nil
}

func getSessionByName(name string) (bool, tmux.Session) {
	server := new(tmux.Server)

	sessions, err := server.ListSessions()
	if err != nil {
		fmt.Println("Error listing sessions:", err)
		return false, tmux.Session{}
	}

	for _, session := range sessions {
		if session.Name == name {
			return true, session
		}
	}

	return false, tmux.Session{}
}
