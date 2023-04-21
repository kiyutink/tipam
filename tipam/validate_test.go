package tipam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOnSubs(t *testing.T) {
	claim := MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false)
	subsPassing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/12", []string{"test", "test_inner"}, false),
	}
	subsFailing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/12", []string{"production"}, false),
	}

	err := ValidateOnSubs(claim, subsPassing)
	assert.Nil(t, err)

	err = ValidateOnSubs(claim, subsFailing)
	assert.NotNil(t, err)
}

func TestValidateOnSupers(t *testing.T) {
	claim := MustParseClaimFromCIDR("10.0.0.0/12", []string{"test", "test_inner"}, false)
	supersPassing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false),
	}
	supersFailing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/8", []string{"production"}, false),
	}

	err := ValidateOnSupers(claim, supersPassing)
	assert.Nil(t, err)

	err = ValidateOnSupers(claim, supersFailing)
	assert.NotNil(t, err)
}
