package main

import (
	"fmt"
	"os"

	gotmux "github.com/GianlucaP106/gotmux/gotmux"
)

func loadSession(project Project, config *Config) error {
	tmux := getTmux()

	session, err := createSession(project, config)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return err
	}

	_, found := os.LookupEnv("TMUX")
	if found {
		tmux.SwitchClient(&gotmux.SwitchClientOptions{TargetSession: session.Name})
	} else {
		session.Attach()
	}

	if config.Active != "" {
		fmt.Println("Active window:", config.Active)
		window, err := session.GetWindowByName(config.Active)
		if err != nil {
			fmt.Println("There is no window with the name: ", config.Active)
			fmt.Println("Please consider changing the active property in your config file", err)
			return err
		}
		window.Select()
	}

	return nil
}

func createSession(project Project, config *Config) (*gotmux.Session, error) {
	tmux := getTmux()

	startWindow := config.Windows[0]

	session, err := tmux.Session(project.Name)
	if session == nil {
		session, err = tmux.NewSession(&gotmux.SessionOptions{
			Name:           project.Name,
			StartDirectory: string(Path(project.Path).Join(startWindow.Path)),
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

			panes, err := windows[0].ListPanes()
			if err != nil {
				fmt.Println("Error listing panes:", err)
				return nil, err
			}

			for _, cmd := range windowConfig.Cmd {
				panes[0].SendKeys(cmd + "\n")
			}

			continue
		}

		window, err := session.NewWindow(&gotmux.NewWindowOptions{
			WindowName:     windowConfig.Name,
			StartDirectory: string(Path(project.Path).Join(windowConfig.Path)),
			DoNotAttach:    true,
		})
		if err != nil {
			fmt.Println("Error creating window:", err)
			return nil, err
		}

		if len(windowConfig.Cmd) == 0 {
			continue
		}

		panes, err := window.ListPanes()
		if err != nil {
			fmt.Println("Error listing panes:", err)
			return nil, err
		}

		for _, cmd := range windowConfig.Cmd {
			panes[0].SendKeys(cmd + "\n")
		}
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
