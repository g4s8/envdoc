package ast

import (
	"errors"
	"fmt"
)

var (
	ErrAstParse   = errors.New("ast parse error")
	ErrFieldParse = fmt.Errorf("parse field: %w", ErrAstParse)
)
