package helper

import (
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/go-uuid"
	"strconv"
)

func ExistsIn(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func FindRule(rules []iaas.RuleDetails, rule iaas.SecurityGroupRuleOpts) *iaas.RuleDetails {
	for _, r := range rules {
		portFrom, _ := strconv.Atoi(rule.PortFrom)
		portTo, _ := strconv.Atoi(rule.PortTo)
		if (ExistsIn(rule.IPs, r.IP) || ExistsIn(rule.IPs, "any")) &&
			r.Direction == rule.Direction &&
			r.Protocol == rule.Protocol &&
			r.PortStart == portFrom &&
			r.PortEnd == portTo {
			return &r
		}
	}
	return nil
}

func GenUUID() (u string) {
	u, _ = uuid.GenerateUUID()
	return u
}
