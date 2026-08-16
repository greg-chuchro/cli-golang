package completion

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/common"
)

// DotnetSuggestHandler handles the dotnet-suggest completion protocol.
type DotnetSuggestHandler struct {
	CompletionHandler ICompletionHandler
}

func NewDotnetSuggestHandler(completionHandler ICompletionHandler) *DotnetSuggestHandler {
	return &DotnetSuggestHandler{CompletionHandler: completionHandler}
}

func (this *DotnetSuggestHandler) HandleDotnetSuggest(context common.ICliContext, args []string, target any) {
	if len(args) < 2 {
		return
	}
	directive := args[0]
	if !strings.HasPrefix(directive, "[suggest") || !strings.HasSuffix(directive, "]") {
		return
	}

	cursorPosition := -1
	if strings.HasPrefix(directive, "[suggest:") {
		posStr := directive[9 : len(directive)-1]
		if n, err := strconv.Atoi(posStr); err == nil {
			cursorPosition = n
		}
	}

	fullCommandLine := args[1]
	commandLine := fullCommandLine
	if cursorPosition >= 0 && cursorPosition < len(fullCommandLine) {
		commandLine = fullCommandLine[:cursorPosition]
	}

	tokens := strings.Split(commandLine, " ")
	if len(tokens) > 0 {
		exeName := filepath.Base(os.Args[0])
		if strings.EqualFold(tokens[0], exeName) {
			tokens = tokens[1:]
		}
	}
	if cursorPosition > len(fullCommandLine) {
		tokens = append(tokens, "")
	}

	completeArgs := append([]string{"__complete"}, tokens...)
	this.CompletionHandler.HandleComplete(context, completeArgs, target)
}

func (this *DotnetSuggestHandler) Register(commandPath string) {
	if commandPath == "" {
		commandPath = os.Args[0]
	}
	if commandPath == "" {
		return
	}
	cmd := exec.Command("dotnet-suggest", "register", "--command-path", commandPath)
	_ = cmd.Run()
}

func (this *DotnetSuggestHandler) Unregister(commandPath string) {
	if commandPath == "" {
		commandPath = os.Args[0]
	}
	if commandPath == "" {
		return
	}
	cmd := exec.Command("dotnet-suggest", "unregister", "--command-path", commandPath)
	_ = cmd.Run()
}
