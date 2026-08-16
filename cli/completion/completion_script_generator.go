package completion

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/greg-chuchro/cli-golang/cli/exceptions"
)

const bashTemplate = `# Bash completion script for __PROGRAM_NAME__
___PROGRAM_NAME___completion() {
    local cur prev words cword
    _init_completion || return
    local completions
    completions=$(__PROGRAM_NAME__ __complete "${words[@]:1}" "$cur" 2>/dev/null)
    COMPREPLY=( $(compgen -W "$completions" -- "$cur") )
}
complete -F ___PROGRAM_NAME___completion __PROGRAM_NAME__`

const zshTemplate = `#compdef __PROGRAM_NAME__
___PROGRAM_NAME___completion() {
    local -a completions
    local -a words
    words=(${(z)LBUFFER})
    local output
    output=$(__PROGRAM_NAME__ __complete "${words[@]:1}" "" 2>/dev/null)
    completions=(${(f)output})
    _describe '__PROGRAM_NAME__ commands' completions
}
___PROGRAM_NAME___completion "$@"`

const powershellTemplate = `# PowerShell completion script for __PROGRAM_NAME__
Register-ArgumentCompleter -Native -CommandName __PROGRAM_NAME__ -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $words = $commandAst.ToString() -split '\s+' | Select-Object -Skip 1
    $completions = & __PROGRAM_NAME__ __complete @words $wordToComplete 2>$null
    $completions | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}`

const fishTemplate = `# Fish completion script for __PROGRAM_NAME__
function ____PROGRAM_NAME___completion
    set -l tokens (commandline -opc)
    set -l args (string join ' ' $tokens[2..-1])
    set -l current (commandline -ct)
    __PROGRAM_NAME__ __complete $tokens[2..-1] "$current" 2>/dev/null
end
complete -c __PROGRAM_NAME__ -f -a "(____PROGRAM_NAME___completion)"`

type shellConfig struct {
	template       string
	getInstallPath func(programName string) string
}

// CompletionScriptGenerator generates and installs shell completion scripts.
type CompletionScriptGenerator struct{}

func NewCompletionScriptGenerator() *CompletionScriptGenerator {
	return &CompletionScriptGenerator{}
}

func (this *CompletionScriptGenerator) shellConfigs() map[string]shellConfig {
	return map[string]shellConfig{
		"bash":       {bashTemplate, func(p string) string { return "/etc/bash_completion.d/" + p }},
		"zsh":        {zshTemplate, func(p string) string { return "/usr/local/share/zsh/site-functions/_" + p }},
		"powershell": {powershellTemplate, func(p string) string { return "completions/" + p + ".ps1" }},
		"pwsh":       {powershellTemplate, func(p string) string { return "completions/" + p + ".ps1" }},
		"fish": {fishTemplate, func(p string) string {
			return filepath.Join(os.Getenv("HOME"), ".config", "fish", "completions", p+".fish")
		}},
	}
}

func (this *CompletionScriptGenerator) SupportedShells() []string {
	configs := this.shellConfigs()
	out := make([]string, 0, len(configs))
	for k := range configs {
		out = append(out, k)
	}
	return out
}

func (this *CompletionScriptGenerator) GenerateScript(shell string, programName string) string {
	config, ok := this.shellConfigs()[strings.ToLower(shell)]
	if !ok {
		panic(exceptions.CliErrors.UnsupportedShell(shell))
	}
	return strings.ReplaceAll(config.template, "__PROGRAM_NAME__", programName)
}

func (this *CompletionScriptGenerator) GetInstallPath(shell string, programName string) string {
	config, ok := this.shellConfigs()[strings.ToLower(shell)]
	if !ok {
		panic(exceptions.CliErrors.UnsupportedShell(shell))
	}
	return config.getInstallPath(programName)
}

func (this *CompletionScriptGenerator) InstallScript(shell string, programName string) {
	script := this.GenerateScript(shell, programName)
	installPath := this.GetInstallPath(shell, programName)
	if dir := filepath.Dir(installPath); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}
	_ = os.WriteFile(installPath, []byte(script), 0o644)
}

func (this *CompletionScriptGenerator) UninstallScript(shell string, programName string) bool {
	installPath := this.GetInstallPath(shell, programName)
	if _, err := os.Stat(installPath); err == nil {
		_ = os.Remove(installPath)
		return true
	}
	return false
}
