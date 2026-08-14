package model

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ultraviolet"
	"github.com/richavery/donk-cli/internal/config"
	"github.com/richavery/donk-cli/internal/home"
	"github.com/richavery/donk-cli/internal/ui/common"
	"github.com/richavery/donk-cli/internal/ui/dialog"
	"github.com/richavery/donk-cli/internal/ui/util"
)

type onboardingStep int

const (
	onboardingStepWelcome onboardingStep = iota
	onboardingStepDependencies
	onboardingStepProvider
	onboardingStepModel
	onboardingStepProject
	onboardingStepSkills
	onboardingStepComplete
)

type onboardingModel struct {
	com   *common.Common
	step  onboardingStep
	done  bool
	error string
	info  string

	hasGit    bool
	hasOllama bool

	providers []string

	skillsSynced bool
	skillsCount  int
}

var _ dialog.Dialog = (*onboardingModel)(nil)

func newOnboardingModel(com *common.Common) *onboardingModel {
	return &onboardingModel{
		com:  com,
		step: onboardingStepWelcome,
	}
}

func (o *onboardingModel) ID() string {
	return "onboarding"
}

func (o *onboardingModel) Init() tea.Cmd {
	return tea.Batch(
		o.checkDependencies(),
		o.detectProviders(),
	)
}

func (o *onboardingModel) checkDependencies() tea.Cmd {
	return func() tea.Msg {
		gitPath, gitErr := lookPath("git")
		ollamaPath, ollamaErr := lookPath("ollama")
		return dependencyCheckMsg{
			hasGit:    gitErr == nil && gitPath != "",
			hasOllama: ollamaErr == nil && ollamaPath != "",
		}
	}
}

func (o *onboardingModel) detectProviders() tea.Cmd {
	return func() tea.Msg {
		providers, err := config.Providers(o.com.Config())
		ids := make([]string, 0)
		if err == nil {
			for _, p := range providers {
				ids = append(ids, string(p.ID))
			}
		}
		return providerDetectMsg{providers: ids, err: err}
	}
}

func (o *onboardingModel) syncSkills() tea.Cmd {
	return func() tea.Msg {
		count := 0
		for _, p := range config.GlobalSkillsDirs() {
			entries, err := os.ReadDir(p)
			if err != nil {
				continue
			}
			for _, e := range entries {
				if !strings.HasPrefix(e.Name(), ".") {
					count++
				}
			}
		}
		return skillsSyncedMsg{count: count}
	}
}

type dependencyCheckMsg struct {
	hasGit    bool
	hasOllama bool
}

type providerDetectMsg struct {
	providers []string
	err       error
}

type skillsSyncedMsg struct {
	count int
}

func (o *onboardingModel) HandleMsg(msg tea.Msg) dialog.Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return o.handleKey(msg)
	case dependencyCheckMsg:
		o.hasGit = msg.hasGit
		o.hasOllama = msg.hasOllama
		o.nextStep()
		return nil
	case providerDetectMsg:
		o.providers = msg.providers
		o.nextStep()
		return nil
	case skillsSyncedMsg:
		o.skillsCount = msg.count
		o.skillsSynced = true
		o.nextStep()
		return nil
	case util.InfoMsg:
		if msg.Type == util.InfoTypeError {
			o.error = msg.Msg
		} else {
			o.info = msg.Msg
		}
		return nil
	}
	return nil
}

func (o *onboardingModel) handleKey(msg tea.KeyPressMsg) dialog.Action {
	switch msg.String() {
	case "esc":
		return dialog.ActionSkipOnboarding{}
	case "enter":
		switch o.step {
		case onboardingStepComplete:
			return dialog.ActionClose{}
		default:
			o.nextStep()
			return nil
		}
	}
	return nil
}

func (o *onboardingModel) nextStep() {
	switch o.step {
	case onboardingStepWelcome:
		o.step = onboardingStepDependencies
	case onboardingStepDependencies:
		o.step = onboardingStepProvider
	case onboardingStepProvider:
		o.step = onboardingStepModel
	case onboardingStepModel:
		o.step = onboardingStepProject
	case onboardingStepProject:
		o.step = onboardingStepSkills
	case onboardingStepSkills:
		o.step = onboardingStepComplete
	case onboardingStepComplete:
		o.done = true
	}
}

func (o *onboardingModel) Draw(scr ultraviolet.Screen, area ultraviolet.Rectangle) *tea.Cursor {
	view := o.render()
	cur := &tea.Cursor{}
	styled := ultraviolet.NewStyledString(view)
	styled.Draw(scr, area)
	return cur
}

func (o *onboardingModel) render() string {
	var b strings.Builder
	b.WriteString(o.com.Styles.Initialize.Header.Render("Welcome to DONK"))
	b.WriteString("\n\n")

	switch o.step {
	case onboardingStepWelcome:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Let's get you set up."))
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render("This will take about a minute."))
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepDependencies:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Checking dependencies..."))
		b.WriteString("\n")
		if o.hasGit {
			b.WriteString(o.com.Styles.Initialize.Content.Render("• Git detected"))
		} else {
			b.WriteString(o.com.Styles.Initialize.Content.Render("• Git not found"))
		}
		b.WriteString("\n")
		if o.hasOllama {
			b.WriteString(o.com.Styles.Initialize.Content.Render("• Ollama detected"))
		} else {
			b.WriteString(o.com.Styles.Initialize.Content.Render("• Ollama not found"))
			b.WriteString("\n")
			b.WriteString(o.com.Styles.Initialize.Content.Render("  Install from https://ollama.com"))
		}
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepProvider:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Providers:"))
		b.WriteString("\n")
		if len(o.providers) == 0 {
			b.WriteString(o.com.Styles.Initialize.Content.Render("No providers configured yet."))
			b.WriteString("\n")
			b.WriteString(o.com.Styles.Initialize.Content.Render("Use /login or configure a provider in your config."))
		} else {
			for _, p := range o.providers {
				b.WriteString(o.com.Styles.Initialize.Content.Render("• " + p))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepModel:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Model setup:"))
		b.WriteString("\n")
		if o.hasOllama {
			b.WriteString(o.com.Styles.Initialize.Content.Render("Ollama is available. Open /models to select a local model."))
		} else {
			b.WriteString(o.com.Styles.Initialize.Content.Render("Add a cloud provider to select a model."))
		}
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepProject:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Project initialization:"))
		b.WriteString("\n")
		cwd := home.Short(o.com.Workspace.WorkingDir())
		b.WriteString(o.com.Styles.Initialize.Content.Render("Current directory: " + cwd))
		b.WriteString("\n")
		needsInit, _ := o.com.Workspace.ProjectNeedsInitialization()
		if needsInit {
			b.WriteString(o.com.Styles.Initialize.Content.Render("This project has not been initialized."))
			b.WriteString("\n")
			b.WriteString(o.com.Styles.Initialize.Content.Render("Press Ctrl+P in the app to initialize it."))
		} else {
			b.WriteString(o.com.Styles.Initialize.Content.Render("Project already initialized."))
		}
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepSkills:
		b.WriteString(o.com.Styles.Initialize.Content.Render("Skills:"))
		b.WriteString("\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render(fmt.Sprintf("Installed skills: %d", o.skillsCount)))
		b.WriteString("\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render("Run install-master-skills.sh to sync the catalog."))
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to continue"))

	case onboardingStepComplete:
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Setup complete!"))
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render("You're ready to use DONK."))
		b.WriteString("\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render("Open /models to choose a model, then start chatting."))
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Accent.Render("Press Enter to launch DONK"))
	}

	if o.error != "" {
		b.WriteString("\n\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render("Error: " + o.error))
	}
	if o.info != "" {
		b.WriteString("\n")
		b.WriteString(o.com.Styles.Initialize.Content.Render(o.info))
	}

	b.WriteString("\n\n")
	b.WriteString(o.com.Styles.Initialize.Content.Render("Press Esc to skip setup"))

	return b.String()
}

// lookPath is a tiny local wrapper around exec.LookPath to avoid importing
// os/exec in the UI layer if that package is unavailable.
func lookPath(name string) (string, error) {
	for _, dir := range []string{"/usr/local/bin", "/usr/bin", "/bin", "/opt/homebrew/bin"} {
		p := dir + "/" + name
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("%s not found", name)
}
