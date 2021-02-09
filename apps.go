package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var mainApp = &cli.App{
	Name:                 "gitflow-toolkit",
	Usage:                "Git Flow ToolKit",
	Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
	Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
	Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
	EnableBashCompletion: true,
	Action: func(c *cli.Context) error {
		return cli.ShowAppHelp(c)
	},
	Commands: []*cli.Command{
		installCmd(),
		uninstallCmd(),
	},
}

var subApps = []*cli.App{
	newBranchApp(FEAT),
	newBranchApp(FIX),
	newBranchApp(DOCS),
	//newBranchApp(STYLE),
	newBranchApp(REFACTOR),
	newBranchApp(CLEAN),
	newBranchApp(TEST),
	//newBranchApp(CHORE),
	//newBranchApp(PERF),
	//newBranchApp(HOTFIX),
	commitApp(),
	checkMessageApp(),
	pushApp(),
}

func newBranchApp(ct CommitType) *cli.App {
	return &cli.App{
		Name:                 "git-" + string(ct),
		Usage:                fmt.Sprintf("Create %s branch", ct),
		UsageText:            fmt.Sprintf("git %s BRANCH", ct),
		Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			err := createBranch(fmt.Sprintf("%s/%s", ct, c.Args().First()))
			if err != nil {
				return fmt.Errorf("failed to create branch %s/%s: %s", ct, c.Args().First(), err)
			}
			return nil
		},
	}
}

func commitApp() *cli.App {
	return &cli.App{
		Name:                 "git-ci",
		Usage:                "Interactive commit",
		UsageText:            "git ci",
		Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return commit()
		},
	}
}

func sendApp() *cli.App {
	return &cli.App{
		Name:                 "git-send",
		Usage:                "Interactive send FeiShu",
		UsageText:            "git send",
		Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return send()
		},
	}
}

func checkMessageApp() *cli.App {
	return &cli.App{
		Name:                 "commit-msg",
		Usage:                "Commit message hook",
		UsageText:            "commit-msg FILE",
		Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			return commitMessageCheck(c.Args().First())
		},
	}
}

func pushApp() *cli.App {
	return &cli.App{
		Name:                 "git-ps",
		Usage:                "Push local branch to remote",
		UsageText:            "git ps",
		Version:              fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors:              []*cli.Author{{Name: "mritd", Email: "mritd@linux.com"}},
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.ShowAppHelp(c)
			}
			return push()
		},
	}
}
