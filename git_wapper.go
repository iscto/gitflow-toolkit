package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/mritd/bubbles/common"

	"github.com/mritd/bubbles/prompt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mritd/bubbles/selector"
)

type MessageType struct {
	Type          CommitType
	ZHDescription string
	ENDescription string
}

type CommitMessage struct {
	Type   CommitType
	Scope  string
	Title  string
	Body   string
	Footer string
	Sob    string
}

func createBranch(name string) error {
	err := repoCheck()
	if err != nil {
		return fmt.Errorf("the current directory is not a git repository")
	}
	return gitCommand(os.Stdout, []string{"checkout", "-b", name})
}

func hasStagedFiles() (bool, error) {
	var buf bytes.Buffer
	err := gitCommand(&buf, []string{"diff", "--cached", "--name-only"})
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(buf.String()) != "", nil
}

func currentBranch() (string, error) {
	var buf bytes.Buffer
	err := gitCommand(&buf, []string{"symbolic-ref", "--short", "HEAD"})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func push() error {
	err := repoCheck()
	if err != nil {
		return fmt.Errorf("the current directory is not a git repository")
	}

	branch, err := currentBranch()
	if err != nil {
		return err
	}
	return gitCommand(os.Stdout, []string{"push", "origin", branch})
}

func commitMessageCheck(f string) error {
	reg := regexp.MustCompile(commitMessagePattern)
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	//commitMsg := string(bs)
	//fmt.Printf(commitMsg)
	commitMsg := reg.FindStringSubmatch(string(bs))
	//i := len(commitMsg)
	//strings :=strconv.Itoa(i)
	//fmt.Printf(strings)

	if len(commitMsg) != 3 {
		return fmt.Errorf(commitMessageCheckFailedTpl)
	}

	return nil
}

func commit() error {
	err := repoCheck()
	if err != nil {
		return fmt.Errorf("the current directory is not a git repository")
	}

	ok, err := hasStagedFiles()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("\nPlease execute the `git add` command to add files before commit.\n")
	}

	cmType, err := commitType()
	if err != nil {
		return err
	}
	//cmScope, err := commitScope()
	//if err != nil {
	//	return err
	//}
	cmTitle, err := commitSubject()
	if err != nil {
		return err
	}
	cmBody, err := commitBody()
	if err != nil {
		return err
	}
	cmFooter, err := commitFooter()
	if err != nil {
		return err
	}
	cmSOB, err := createSOB()
	if err != nil {
		return err
	}

	msg := CommitMessage{
		Type:    cmType.Type,
		//Scope:   cmScope,
		Title:  cmTitle,
		Body:   cmBody,
		Footer: cmFooter,
		Sob:    cmSOB,
	}
	if msg.Body == "" {
		msg.Body = cmTitle
	}

	f, err := ioutil.TempFile("", "git-commit")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	tpl, _ := template.New("").Parse(commitMessageTpl)
	err = tpl.Execute(f, msg)
	if err != nil {
		return err
	}

	err = gitCommand(os.Stdout, []string{"commit", "-F", f.Name()})
	if err != nil {
		return err
	}

	fmt.Println("\n" + common.FontColor(commitWarn, "2"))

	return nil
}

func send() error {

	cmTitle, err := sendTitle()
	if err != nil {
		return err
	}
	cmBody, err := sendBody()
	if err != nil {
		return err
	}
	cmSOB, err := creator()
	if err != nil {
		return err
	}

	msg := CommitMessage{
		Title:  cmTitle,
		Body:   cmBody,
		Sob:    cmSOB,
	}
	if msg.Body == "" {
		msg.Body = cmTitle
	}

	f, err := ioutil.TempFile("", "send")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	tpl, _ := template.New("").Parse(sendMessageTpl)
	err = tpl.Execute(f, msg)
	if err != nil {
		return err
	}

	//err = gitCommand(os.Stdout, []string{"commit", "-F", f.Name()})
	//if err != nil {
	//	return err
	//}

	fmt.Println("\n" + common.FontColor(f.Name(), "2"))

	return nil
}

func commitType() (MessageType, error) {
	m := &selector.Model{
		Data: []interface{}{
			MessageType{Type: FEAT, ZHDescription: "实现 Task 功能", ENDescription: "Introducing new features"},
			MessageType{Type: FIX, ZHDescription: "修复 Bug", ENDescription: "Bug fix"},
			MessageType{Type: DOCS, ZHDescription: "补充注释和文档", ENDescription: "Writing docs"},
			//MessageType{Type: STYLE, ZHDescription: "调整格式", ENDescription: "Improving structure/format of the code"},
			MessageType{Type: TEST, ZHDescription: "增加测试代码", ENDescription: "When adding missing tests"},
			MessageType{Type: REFACTOR, ZHDescription: "重构代码", ENDescription: "Refactoring code"},
			MessageType{Type: CLEAN, ZHDescription: "重构代码", ENDescription: "Refactoring code"},
			//MessageType{Type: CHORE, ZHDescription: "CI/CD 变动", ENDescription: "Changing CI/CD"},
			//MessageType{Type: PERF, ZHDescription: "性能优化", ENDescription: "Improving performance"},
			//MessageType{Type: HOTFIX, ZHDescription: "紧急修复", ENDescription: "Bug fix urgently"},
		},
		PerPage: 6,
		// Use the arrow keys to navigate: ↓ ↑ → ←
		// Select Commit Type:

		//HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select Commit Type:"),
		HeaderFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			headerTpl := `Use the arrow keys to navigate: ↓ ↑ → ←
Select Commit Type:`
			return common.FontColor(fmt.Sprintf(headerTpl), "222")
		},
		// [1] feat (Introducing new features)
		SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(MessageType)
			return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Type, t.ZHDescription), "220")
		},
		// 2. fix (Bug fix)
		UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(MessageType)
			return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Type, t.ZHDescription), "14")
		},
		// --------- Commit Type ----------
		// Type: feat
		// Description: 新功能(Introducing new features)
		FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := m.PageSelected().(MessageType)
			footerTpl := `--------- Commit Type ----------
Type: %s
Description: %s(%s)`
			return common.FontColor(fmt.Sprintf(footerTpl, t.Type, t.ZHDescription, t.ENDescription), "222")
		},
		FinishedFunc: func(s interface{}) string {
			mt := s.(MessageType)
			return common.FontColor("✔ Type: ", selector.ColorFinished) + string(mt.Type) + "\n"
		},
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return MessageType{}, err
	}

	if m.Canceled() {
		return MessageType{}, fmt.Errorf("user has cancelled this commit")
	}

	return m.Selected().(MessageType), nil
}

func commitScope() (string, error) {
	m := &prompt.Model{
		Prompt:       common.FontColor("Scope: ", "2"),
		ValidateFunc: prompt.VFNotBlank,
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}
	return m.Value(), nil
}

func sendTitle() (string, error) {
	m := &prompt.Model{
		Prompt:       common.FontColor("Title: ", "2"),
		ValidateFunc: prompt.VFNotBlank,
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}
	return m.Value(), nil
}

func commitSubject() (string, error) {
	m := &prompt.Model{
		Prompt:       common.FontColor("Title: ", "2"),
		ValidateFunc: prompt.VFNotBlank,
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}
	return m.Value(), nil
}

func commitBody() (string, error) {
	m := &prompt.Model{
		Prompt: common.FontColor("Body: ", "2"),
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	value := m.Value()
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}

	reg := regexp.MustCompile(commitBodyEditPattern)
	if reg.MatchString(value) {
		return openEditor()
	}
	return m.Value(), nil
}

func sendBody() (string, error) {
	m := &prompt.Model{
		Prompt: common.FontColor("Body: ", "2"),
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	value := m.Value()
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}

	reg := regexp.MustCompile(commitBodyEditPattern)
	if reg.MatchString(value) {
		return openEditor()
	}
	return m.Value(), nil
}

func commitFooter() (string, error) {
	m := &prompt.Model{
		Prompt: common.FontColor("Footer: ", "2"),
	}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		return "", err
	}
	if m.Canceled() {
		return "", fmt.Errorf("user has cancelled this commit")
	}
	return m.Value(), nil
}

func createSOB() (string, error) {
	name, email, err := gitAuthor()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Signed-off-by: %s %s", name, email), nil
}

func creator() (string, error) {
	name, _, err := gitAuthor()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("发布人 %s ", name), nil
}

func gitAuthor() (string, string, error) {
	name := "Undefined"
	email := "Undefined"

	var buf bytes.Buffer
	err := gitCommand(&buf, []string{"var", "GIT_AUTHOR_IDENT"})
	if err != nil {
		return "", "", err
	}

	authorInfo := strings.Fields(buf.String())
	if len(authorInfo) > 1 && authorInfo[0] != "" {
		name = authorInfo[0]
	}
	if len(authorInfo) > 2 && authorInfo[1] != "" {
		email = authorInfo[1]
	}
	return name, email, nil
}

func repoCheck() error {
	return gitCommand(ioutil.Discard, []string{"rev-parse", "--show-toplevel"})
}
