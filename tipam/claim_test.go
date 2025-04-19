package tipam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSubs(t *testing.T) {
	claim := MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false)
	subsPassing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/12", []string{"test", "test_inner"}, false),
	}
	subsFailing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/12", []string{"production"}, false),
	}

	err := claim.ValidateSubs(subsPassing)
	assert.Nil(t, err)

	err = claim.ValidateSubs(subsFailing)
	assert.NotNil(t, err)
}

func TestValidateSupers(t *testing.T) {
	claim := MustParseClaimFromCIDR("10.0.0.0/12", []string{"test", "test_inner"}, false)
	supersPassing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/8", []string{"test"}, false),
	}
	supersFailing := []*Claim{
		MustParseClaimFromCIDR("10.0.0.0/8", []string{"production"}, false),
	}

	err := claim.ValidateSupers(supersPassing)
	assert.Nil(t, err)

	err = claim.ValidateSupers(supersFailing)
	assert.NotNil(t, err)
}
