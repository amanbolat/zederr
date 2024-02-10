package parser

import (
	"fmt"
	"text/template/parse"

	"github.com/amanbolat/zederr/pkg/codegen/core"
	"github.com/k0kubun/pp/v3"
)

// Parser is an implementation of core.Parser interface.
type Parser struct {
	params     []core.Param
	leftDelim  string
	rightDelim string
	debug      bool
	parseTree  *parse.Tree
}

// Config is a configuration for Parser.
type Config struct {
	Debug bool
}

// NewParser returns a new instance of Parser.
func NewParser(leftDelim, rightDelim string, cfg *Config) *Parser {
	var defCfg Config
	if cfg != nil {
		defCfg = *cfg
	}

	return &Parser{
		params:     nil,
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
		debug:      defCfg.Debug,
	}
}

// Parse implements core.Parser interface.
func (p *Parser) Parse(txt string) (_ map[string]core.Param, _ string, err error) {
	p.reset()

	m, err := parse.Parse("", txt, p.leftDelim, p.rightDelim, map[string]any{
		"string": func() {},
		"int":    func() {},
		"bool":   func() {},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse the text with text/template parser: %w", err)
	}

	tree, ok := m[""]
	if !ok {
		return nil, "", fmt.Errorf("parse.Tree returned by text/template parser is emtpy")
	}

	if p.debug {
		_, _ = pp.Printf("parsed tree:\n%v\n", tree)
	}

	err = p.parse(tree.Root)
	if err != nil {
		return nil, "", err
	}

	fieldMap := make(map[string]core.Param)

	for _, f := range p.params {
		fieldMap[f.Name] = f
	}

	_, localizableText := tree.ErrorContext(tree.Root)

	return fieldMap, localizableText, nil
}

func (p *Parser) parse(root parse.Node) error {
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

		err := p.tryParseParam(node.Pipe, node.Pipe.Cmds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) tryParseParam(pipe *parse.PipeNode, cmdNodes []*parse.CommandNode) error {
	param := core.Param{
		Name: "",
		Type: "",
	}

	if len(cmdNodes[0].Args) == 0 {
		return nil
	}

	paramCmdNode := cmdNodes[0]

	pp.Println(cmdNodes[0].Args)

	// Retrieve param name.
	switch node := paramCmdNode.Args[0].(type) {
	case *parse.FieldNode:
		paramName, err := p.extractParamNameFromFieldNode(node)
		if err != nil {
			return err
		}

		param.Name = paramName
	// Handle the case when the type function is before the param name:
	//
	// {{ string .Param}}
	case *parse.IdentifierNode:
		if len(node.Ident) == 0 {
			return nil
		}
		param.Type = node.Ident

		if len(paramCmdNode.Args) < 2 {
			return fmt.Errorf("field name was not found; pos [%d]", paramCmdNode.Pos)
		}

		fieldNode, ok := paramCmdNode.Args[1].(*parse.FieldNode)
		if !ok {
			return fmt.Errorf("failed to cast to FieldNode; pos [%d]", paramCmdNode.Pos)
		}

		paramName, err := p.extractParamNameFromFieldNode(fieldNode)
		if err != nil {
			return err
		}

		param.Name = paramName
	// If a node type is not handled, just ignore it.
	default:
		return nil
	}

	// Handle the case when the type function is before the field name:
	// Example: {{ string .Param}}
	if len(cmdNodes) >= 2 && len(cmdNodes[1].Args) > 0 && param.Type != "" {
		fieldNode, ok := cmdNodes[1].Args[0].(*parse.FieldNode)
		if !ok {
			return fmt.Errorf("field node was not found in second command node; pos [%d]", cmdNodes[1].Pos)
		}

		paramName, err := p.extractParamNameFromFieldNode(fieldNode)
		if err != nil {
			return err
		}

		param.Name = paramName

		pipe.Cmds = pipe.Cmds[1:]
	}

	// Try to get the type of the param if it's provided after the param name.
	// Example: {{ .Param | string }}
	//
	// All the cases where more than one pipe command is used are ignored.
	// Example: {{ .Param | string | int }}
	//
	// It should run only in case when the type was not found.
	if len(cmdNodes) >= 2 && param.Type == "" {
		if len(cmdNodes[1].Args) == 0 {
			return nil
		}

		typeNode, ok := (cmdNodes[1].Args[0]).(*parse.IdentifierNode)
		if !ok {
			_, c := p.parseTree.ErrorContext(typeNode)
			return fmt.Errorf("failed to cast the expected parse.Node to parse.IdentifierNode to get the field type; param name [%s], context [%s]", param.Name, c)
		}
		param.Type = typeNode.Ident

		// Remove the node with the Param type from the parse tree, as i18n
		// module doesn't have type functions like `string` and `int`.
		pipe.Cmds = pipe.Cmds[:1]
	}

	if param.Type == "" {
		param.Type = "any"
	}

	p.params = append(p.params, param)

	return nil
}

func (p *Parser) extractParamNameFromFieldNode(fieldNode *parse.FieldNode) (string, error) {
	if len(fieldNode.Ident) == 0 {
		return "", fmt.Errorf("field name was not found; parse.FieldNode has no Idents; pos [%d]", fieldNode.Pos)
	}

	if len(fieldNode.Ident) > 1 {
		return "", fmt.Errorf("field name contains more than one identifier; pos [%d]. "+
			"You are probably have a field with name such as `.Param.NestedField`. "+
			"It's not allowed. You shoud rename it to `FieldNestedField`", fieldNode.Pos)
	}

	return fieldNode.Ident[0], nil
}

func (p *Parser) reset() {
	p.params = nil
	p.parseTree = nil
}
