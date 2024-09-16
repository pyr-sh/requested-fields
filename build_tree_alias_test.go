package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTreeAlias(t *testing.T) {
	var graphqlQueryA string = `
	  {
		  search {
		    filters
		  }

		  best: search {
		    connection
		  }
		}
	`

	expectedTreeA := map[string][]string{
		"":       {"search"},
		"search": {"filters", "connection"}}

	generatedTreeA := BuildTree(graphqlQueryA, Variables{})

	assert.Equal(t, expectedTreeA, generatedTreeA)

	var graphqlQueryB string = `
	  {
		  search {
		    filters
		  }

		  best: search {
		    connection
		  }

		  worst:search {
		    term
		  }
		}
	`

	expectedTreeB := map[string][]string{
		"":       {"search", "best", "worst"},
		"best":   {"connection"},
		"search": {"filters"},
		"worst":  {"term"}}

	generatedTreeB := BuildTreeUsingAliases(graphqlQueryB, Variables{})

	assert.Equal(t, expectedTreeB, generatedTreeB)

	var graphqlQueryC string = `
	{
	  user(id: 3) {
	    id
	    custom_name: name
	    birthday
	  }

	  custom_user: user(id: 4) {
	    id
	    name
	    age
	  }
	}
	`
	expectedTreeC := map[string][]string{
		"":            {"user", "custom_user"},
		"user":        {"id", "custom_name", "birthday"},
		"custom_user": {"id", "name", "age"}}

	generatedTreeC := BuildTreeUsingAliases(graphqlQueryC, Variables{})

	assert.Equal(t, expectedTreeC, generatedTreeC)
}
