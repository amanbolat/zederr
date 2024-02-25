package core

import (
	"fmt"
	"text/template/parse"

	"github.com/k0kubun/pp/v3"
)

// TemplateValidator is an implementation of core.TemplateValidator interface.
type TemplateValidator struct {
	arguments  map[string]struct{}
	leftDelim  string
	rightDelim string
	debug      bool
	parseTree  *parse.Tree
}

// TemplateValidatorConfig is a configuration for TemplateValidator.
type TemplateValidatorConfig struct {
	Debug     bool
	Arguments map[string]struct{}
}

// NewTemplateValidator returns a new instance of TemplateValidator.
func NewTemplateValidator(cfg *TemplateValidatorConfig) *TemplateValidator {
	var defCfg TemplateValidatorConfig
	if cfg != nil {
		defCfg = *cfg
	}

	return &TemplateValidator{
		arguments:  defCfg.Arguments,
		leftDelim:  "{{",
		rightDelim: "}}",
		debug:      defCfg.Debug,
		parseTree:  nil,
	}
}

// Validate implements core.TemplateValidator interface.
func (p *TemplateValidator) Validate(txt string) (err error) {
	p.reset()

	m, err := parse.Parse("", txt, p.leftDelim, p.rightDelim)
	if err != nil {
		return fmt.Errorf("failed to parse the text with text/template parser: %w", err)
	}

	tree, ok := m[""]
	if !ok {
		return fmt.Errorf("parse.Tree returned by text/template parser is emtpy")
	}

	if p.debug {
		_, _ = pp.Printf("parsed tree:\n%v\n", tree)
	}

	err = p.parse(tree.Root)
	if err != nil {
		return err
	}

	return nil
}

func (p *TemplateValidator) parse(root parse.Node) error {
	switch node := root.(type) {
	case *parse.ListNode:
		for _, childNode := range node.Nodes {
			err := p.parse(childNode)
			if err != nil {
				return err
			}
		}
	case *parse.ActionNode:
		if node.Pipe == nil {
			return nil
		}

		if len(node.Pipe.Cmds) == 0 {
			return nil
		}

		err := p.tryParseParam(node.Pipe.Cmds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TemplateValidator) tryParseParam(cmdNodes []*parse.CommandNode) error {
	if len(cmdNodes[0].Args) == 0 {
		return nil
	}

	paramCmdNode := cmdNodes[0]

	switch node := paramCmdNode.Args[0].(type) {
	case *parse.FieldNode:
		argName, err := p.extractParamNameFromFieldNode(node)
		if err != nil {
			return err
		}

		err = p.checkArgument(argName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected type of the first argument of the command node; expected *parse.FieldNode; got %T", node)
	}

	return nil
}

func (p *TemplateValidator) extractParamNameFromFieldNode(fieldNode *parse.FieldNode) (string, error) {
	if len(fieldNode.Ident) == 0 {
		return "", fmt.Errorf("field name was not found; parse.FieldNode has no Idents; pos [%d]", fieldNode.Pos)
	}

	if len(fieldNode.Ident) > 1 {
		return "", fmt.Errorf("field name contains more than one identifier; pos [%d]. "+
			"You are probably have a field with name such as `.Argument.NestedField`. "+
			"It's not allowed. You shoud rename it to `FieldNestedField`", fieldNode.Pos)
	}

	return fieldNode.Ident[0], nil
}

func (p *TemplateValidator) reset() {
	p.parseTree = nil
}

func (p *TemplateValidator) checkArgument(argName string) error {
	if _, ok := p.arguments[argName]; !ok {
		return fmt.Errorf("argument with name [%s] was not found in the list of arguments", argName)
	}

	return nil
}
