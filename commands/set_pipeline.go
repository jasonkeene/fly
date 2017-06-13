package commands

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/web"
	"github.com/concourse/fly/commands/internal/flaghelpers"
	"github.com/concourse/fly/commands/internal/setpipelinehelpers"
	"github.com/concourse/fly/rc"
	"github.com/concourse/fly/template"
	"github.com/tedsuo/rata"
)

type SetPipelineCommand struct {
	Pipeline        flaghelpers.PipelineFlag       `short:"p"  long:"pipeline" required:"true"      description:"Pipeline to configure"`
	Config          atc.PathFlag                   `short:"c"  long:"config" required:"true"        description:"Pipeline configuration file"`
	Var             []flaghelpers.VariablePairFlag `short:"v"  long:"var" value-name:"[SECRET=KEY]" description:"Variable flag that can be used for filling in template values in configuration"`
	VarsFrom        []atc.PathFlag                 `short:"l"  long:"load-vars-from"                description:"Variable flag that can be used for filling in template values in configuration from a YAML file"`
	SkipInteractive bool                           `short:"n"  long:"non-interactive"               description:"Skips interactions, uses default values"`
}

func (command *SetPipelineCommand) Execute(args []string) error {
	configPath := command.Config
	templateVariablesFiles := command.VarsFrom
	pipelineName := string(command.Pipeline)

	templateVariables := map[string]interface{}{}
	oldTemplateVariables := template.Variables{}
	for _, v := range command.Var {
		templateVariables[v.Name] = v.Value
		oldTemplateVariables[v.Name] = v.OldValue
	}

	target, err := rc.LoadTarget(Fly.Target)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	webRequestGenerator := rata.NewRequestGenerator(target.Client().URL(), web.Routes)

	atcConfig := setpipelinehelpers.ATCConfig{
		Team:                target.Team(),
		PipelineName:        pipelineName,
		WebRequestGenerator: webRequestGenerator,
		SkipInteraction:     command.SkipInteractive,
	}

	return atcConfig.Set(configPath, templateVariables, oldTemplateVariables, templateVariablesFiles)
}
