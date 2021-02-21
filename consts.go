package main

type CommitType string

const (
	FEAT     CommitType = "feature"
	FIX      CommitType = "bugfix"
	DOCS     CommitType = "doc"
	TEST     CommitType = "test"
	REFACTOR CommitType = "refactor"
	CLEAN 	 CommitType = "clean"
)

const commitWarn = "✔ Follow the arc diff at https://phab.gotokeep.com ask others to help you Accept Revision."

const commitMessagePattern = `^(feature|bugfix|doc|test|refactor|clean):{1}(\s.*)`
//const commitMessagePattern = `^(feat|fix|docs|style|refactor|test|chore|perf|hotfix)\((\S.*)\):\s(\S.*)|^Merge.*`
const commitBodyEditPattern = `^\/\/\s*(?i)edit`

const commitMessageTpl = `{{.Type}}: {{.Title}}

{{ .Body }}

{{ .Footer }}

{{ .Sob }}

`

const sendMessageTpl = `{{.Title}}
{{ .Body }}
{{ .Sob }}
`

//const commitMessageTpl = `{{ .Type }}({{ .Scope }}): {{ .Title }}
//
//{{ .Body }}
//
//{{ .Footer }}
//
//{{ .Sob }}
//`
const commitMessageCheckFailedTpl = `
########################################################
##                                                    ##
##    💔 The commit message is not standardized.      ##
##    💔 It must match the regular expression:        ##
##                                                    ##
## ^(feature|bugfix|doc|test|refactor|clean):{1}(\s.*)##
##                                                    ##
########################################################`