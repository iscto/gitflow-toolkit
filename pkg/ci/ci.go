package ci

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mritd/gitflow-toolkit/pkg/consts"
	"github.com/mritd/gitflow-toolkit/pkg/util"
	"github.com/mritd/promptui"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

type CommitTypeMessage struct {
	Type          consts.CommitType
	ZHDescription string
	ENDescription string
}

type CommitMessage struct {
	Type    consts.CommitType
	Scope   string
	Subject string
	Body    string
	Footer  string
	Sob     string
}

// 检查当前位置是否为 git 项目
func CheckGitProject() bool {
	return exec.Command("git", "rev-parse", "--show-toplevel").Run() == nil
}

// 选择提交类型
func SelectCommitType() consts.CommitType {

	commitTypes := []CommitTypeMessage{
		{Type: consts.FEAT, ZHDescription: "新功能", ENDescription: "Introducing new features"},
		{Type: consts.FIX, ZHDescription: "修复 Bug", ENDescription: "Bug fix"},
		{Type: consts.DOCS, ZHDescription: "添加文档", ENDescription: "Writing docs"},
		{Type: consts.STYLE, ZHDescription: "调整格式", ENDescription: "Improving structure/format of the code"},
		{Type: consts.REFACTOR, ZHDescription: "重构代码", ENDescription: "Refactoring code"},
		{Type: consts.TEST, ZHDescription: "增加测试", ENDescription: "When adding missing tests"},
		{Type: consts.CHORE, ZHDescription: "CI/CD 变动", ENDescription: "Changing CI/CD"},
		{Type: consts.PERF, ZHDescription: "性能优化", ENDescription: "Improving performance"},
		{Type: consts.EXIT, ZHDescription: "退出", ENDescription: "Exit commit"},
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "❯ {{ .Type | cyan }} ({{ .ENDescription | cyan }})",
		Inactive: "  {{ .Type | white }} ({{ .ENDescription | white }})",
		Selected: "{{ \"❯ Type\" | green }}: {{ .Type }}",
		Details: `
--------- Commit Type ----------
{{ "Type:" | faint }}	{{ .Type }}
{{ "Description:" | faint }}	{{ .ZHDescription }}({{ .ENDescription }})`,
	}

	searcher := func(input string, index int) bool {
		commitType := commitTypes[index]
		cmType := strings.Replace(strings.ToLower(string(commitType.Type)), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(cmType, input)
	}

	prompt := promptui.Select{
		Label:     "Select Commit Type:",
		Items:     commitTypes,
		Templates: templates,
		Size:      9,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	util.CheckAndExit(err)

	return commitTypes[i].Type
}

// 输入影响范围
func InputScope() string {

	validate := func(input string) error {
		reg := regexp.MustCompile("\\s+")
		if input == "" || reg.ReplaceAllString(input, "") == "" {
			return errors.New("scope is blank")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "❯ Scope:",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	util.CheckAndExit(err)
	return result

}

// 输入提交主题
func InputSubject() string {

	validate := func(input string) error {
		reg := regexp.MustCompile("\\s+")
		if input == "" || reg.ReplaceAllString(input, "") == "" {
			return errors.New("subject is blank")
		}
		if r := []rune(input); len(r) > 50 {
			return errors.New("subject too long")
		}

		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "❯ Subject:",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	util.CheckAndExit(err)
	return result
}

// 输入完整提交信息
func InputBody() string {

	validate := func(input string) error {
		reg := regexp.MustCompile("\\s+")
		if input == "" || reg.ReplaceAllString(input, "") == "" {
			return errors.New("body is blank")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "❯ Body:",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	util.CheckAndExit(err)
	return result
}

// 输入超长文本信息
func InputBigBody() string {

	f, err := ioutil.TempFile("", "gitflow-toolkit")
	util.CheckAndExit(err)
	defer os.Remove(f.Name())

	// write utf8 bom
	bom := []byte{0xef, 0xbb, 0xbf}
	_, err = f.Write(bom)
	util.CheckAndExit(err)

	// 获取系统编辑器
	editor := "vim"
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}

	// 执行编辑文件
	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	util.CheckAndExit(cmd.Run())
	raw, err := ioutil.ReadFile(f.Name())
	util.CheckAndExit(err)
	body := string(bytes.TrimPrefix(raw, bom))

	return body
}

// 输入提交关联信息
func InputFooter() string {

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "❯ Footer:",
		Templates: templates,
	}

	result, err := prompt.Run()
	util.CheckAndExit(err)
	return result
}

// 生成 SOB 签名
func GenSOB() string {
	cmd := exec.Command("git", "var", "GIT_AUTHOR_IDENT")
	buf, err := cmd.CombinedOutput()
	util.CheckAndExit(err)
	if string(buf) == "" {
		return ""
	}

	author := "Undefined"
	email := "Undefined"
	authorInfo := strings.Fields(string(buf))

	if authorInfo[0] != "" {
		author = authorInfo[0]
	}
	if authorInfo[1] != "" {
		email = authorInfo[1]
	}

	return "Signed-off-by: " + author + " " + email
}

// 提交
func Commit(cm *CommitMessage) {
	t, err := template.New("commitMessage").Parse(consts.CommitTpl)
	util.CheckAndExit(err)
	f, err := ioutil.TempFile("", "git-commit")
	defer os.Remove(f.Name())
	util.CheckAndExit(err)
	t.Execute(f, cm)
	cmd := exec.Command("git", "commit", "-F", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	util.CheckAndExit(cmd.Run())
	fmt.Println("\nAlways code as if the guy who ends up maintaining your code will be a violent psychopath who knows where you live.")
}