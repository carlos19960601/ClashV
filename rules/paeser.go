package rules

import (
	"fmt"

	C "github.com/carlos19960601/ClashV/constant"
	RC "github.com/carlos19960601/ClashV/rules/common"
)

func ParseRule(tp, payload, target string, params []string, subRules map[string][]C.Rule) (parsed C.Rule, parseErr error) {
	switch tp {
	case "DOMAIN":
		parsed = RC.NewDomain(payload, target)
	case "DOMAIN-SUFFIX":
		parsed = RC.NewDomainSuffix(payload, target)
	default:
		parseErr = fmt.Errorf("不支持的rule")
	}

	if parseErr != nil {
		return nil, parseErr
	}

	return
}
