package main

import (
	"fmt"
	gotmux "github.com/GianlucaP106/gotmux/gotmux"
)

func loadSession(name string, config *Config) error {
	tmux := getTmux()

	session, err := createSession(name, config)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return err
	}

	tmux.SwitchClient(&gotmux.SwitchClientOptions{TargetSession: session.Name})

	return nil
}

func createSession(name string, config *Config) (*gotmux.Session, error) {
	tmux := getTmux()

	startWindow := config.Windows[0]

	session, err := tmux.Session(name)
	if session == nil {
		session, err = tmux.NewSession(&gotmux.SessionOptions{
			Name:           PROJECT_NAME,
			StartDirectory: fmt.Sprintf("%s/%s/%s", DIRPATH, PROJECT_NAME, startWindow.Path),
			ShellCommand:   startWindow.Cmd,
		})

		if err != nil {
			fmt.Println("Error creating session:", err)
			return nil, err
		}
	} else {
		return session, nil
	}

	for idx, windowConfig := range config.Windows {
		if idx == 0 {
			windows, err := session.ListWindows()
			if err != nil {
				fmt.Println("Error getting window:", err)
				return nil, err
			}

			windows[0].Rename(windowConfig.Name)

			continue
		}

		window, err := session.NewWindow(&gotmux.NewWindowOptions{
			WindowName:     windowConfig.Name,
			StartDirectory: fmt.Sprintf("%s/%s/%s", DIRPATH, PROJECT_NAME, startWindow.Path),
			DoNotAttach:    true,
		})
		if err != nil {
			fmt.Println("Error creating window:", err)
			return nil, err
		}

		if windowConfig.Cmd == "" {
			continue
		}

		panes, err := window.ListPanes()
		if err != nil {
			fmt.Println("Error listing panes:", err)
			return nil, err
		}

		panes[0].SendKeys(windowConfig.Cmd + "\n")
	}

	return session, nil
}

func getTmux() *gotmux.Tmux {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		panic(err)
	}

	return tmux
}
