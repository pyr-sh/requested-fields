package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Variables map[string]interface{}

func TestBuildTree(t *testing.T) {
	var graphqlQueryA string = `

		query productsSearch(

		  $products_search: ProductsSearchInput!) {

		  search(products: $products_search)


		  {
		    term
		    products {
		      edges {
		        node {
		          id           custom_title: title
		          seller {
		            ...SellerData
		          }
		        }
		        cursor
		      }
		    }
		  }

		  search_users {
		    term
		    users {
		      edges {
		        node {
		          id           custom_title: title
		          seller {...SellerData}
		        }
		        cursor
		      }
		    }
		  }
		}

		fragment SellerData on User {
		  id, ...SellerDataB
		} 

		fragment

		SellerDataB

		on

		  User {
		  name
		} `

	// search {
	//   term
	//   products {
	//     edges {
	//       node {
	//         id
	//         title
	//         seller {
	//           id
	//           name
	//         }
	//       }
	//       cursor
	//     }
	//   }
	// }
	// search_users {
	//   term
	//   users {
	//     edges {
	//       node {
	//         id
	//         title
	//         seller {
	//           id
	//           name
	//         }
	//       }
	//       cursor
	//     }
	//   }
	// }

	var graphqlQueryB string = `
		query {
		  users {
		    id
		    title
		  }
		}`

	var graphqlQueryC string = `
		query {
		  hello
		}`

	var graphqlQueryD string = `
		{
		  user(id: 3) {
		    id
		    name
		  }
		}`

	var graphqlQueryE string = `
		{
		  user(id: 3) {
		    id
		    name
		  }

		  custom_user: user(id: 4) {
		    id
		    name
		    age
		  }
		}`

	var graphqlQueryF string = `{
		  users {
		    users {
		      users {
		        name
		      }
		    }
		  }
		}`

	var graphqlQueryG string = `
		{
		  ...Frag
		}

		fragment Frag on SomeType {
		  field {
		    sub_field
		  }
		}`

	expectedTreeA := map[string][]string{
		"":                                     {"search", "search_users"},
		"search":                               {"term", "products"},
		"search.products":                      {"edges"},
		"search.products.edges":                {"node", "cursor"},
		"search.products.edges.node":           {"id", "title", "seller"},
		"search.products.edges.node.seller":    {"id", "name"},
		"search_users":                         {"term", "users"},
		"search_users.users":                   {"edges"},
		"search_users.users.edges":             {"node", "cursor"},
		"search_users.users.edges.node":        {"id", "title", "seller"},
		"search_users.users.edges.node.seller": {"id", "name"},
	}

	generatedTreeA := BuildTree(graphqlQueryA, Variables{})

	assert.Equal(t, expectedTreeA[""], generatedTreeA[""])

	assert.Equal(t, expectedTreeA, generatedTreeA)

	expectedTreeB := map[string][]string{
		"":      {"users"},
		"users": {"id", "title"},
	}

	generatedTreeB := BuildTree(graphqlQueryB, Variables{})

	assert.Equal(t, expectedTreeB, generatedTreeB)

	expectedTreeC := map[string][]string{
		"": {"hello"},
	}

	generatedTreeC := BuildTree(graphqlQueryC, Variables{})

	assert.Equal(t, expectedTreeC, generatedTreeC)

	expectedTreeD := map[string][]string{
		"":     {"user"},
		"user": {"id", "name"},
	}

	generatedTreeD := BuildTree(graphqlQueryD, Variables{})

	assert.Equal(t, expectedTreeD, generatedTreeD)

	expectedTreeE := map[string][]string{
		"":     {"user"},
		"user": {"id", "name", "age"},
	}

	generatedTreeE := BuildTree(graphqlQueryE, Variables{})

	assert.Equal(t, expectedTreeE, generatedTreeE)

	expectedTreeF := map[string][]string{
		"":                  {"users"},
		"users":             {"users"},
		"users.users":       {"users"},
		"users.users.users": {"name"}}

	generatedTreeF := BuildTree(graphqlQueryF, Variables{})

	assert.Equal(t, expectedTreeF, generatedTreeF)

	expectedTreeG := map[string][]string{
		"":      {"field"},
		"field": {"sub_field"}}

	generatedTreeG := BuildTree(graphqlQueryG, Variables{})

	assert.Equal(t, expectedTreeG, generatedTreeG)

	var graphqlQueryH string = `
		query (
		  $product_id: ID!,
		  $first: Int!
		){
		  product(id: $product_id) {
		    id
		  }

		  search(products: $search) {
		    term
		  }

		  other_a: search(products: $search) {
		    products(
		      first: $first,
		      sort: $sort,
		    ) {
		      total
		    }
		  }

		  other_b: search(products: $search) {
		    products(
		      first: $first,
		      sort: $sort
		    ) {
		      total
		    }
		  }
		}
	`
	expectedTreeH := map[string][]string{
		"":                 {"product", "search", "other_a", "other_b"},
		"other_a":          {"products"},
		"other_a.products": {"total"},
		"other_b":          {"products"},
		"other_b.products": {"total"},
		"product":          {"id"},
		"search":           {"term"}}

	generatedTreeH := BuildTreeUsingAliases(graphqlQueryH, Variables{})

	assert.Equal(t, expectedTreeH, generatedTreeH)

	// -------------------------------------

	var graphqlQueryI string = `
		query ProductsSearchPage($include_aggregations: Boolean!) {
		  search {
		    aggregations {
		      departments {
		        ...departmentAggregationFields
		        __typename
		      }
		    }
		    some_field
		  }
		}

		fragment departmentAggregationFields on DepartmentAggregation {
		  slug
		  name
		  __typename
		}

	`
	expectedTreeI := map[string][]string{
		"":                                {"search"},
		"search":                          {"aggregations", "some_field"},
		"search.aggregations":             {"departments"},
		"search.aggregations.departments": {"slug", "name", "__typename"}}

	generatedTreeI := BuildTreeUsingAliases(graphqlQueryI, Variables{})

	assert.Equal(t, expectedTreeI, generatedTreeI)

	// -------------

	var graphqlQueryJ string = `
		query ProductsSearchPage($include_aggregations: Boolean!) {
		  search {
		    aggregations @include(if: true) {
		      departments {
		        ...departmentAggregationFields
		        __typename
		      }
		    }
		    some_field
		  }
		}

		fragment departmentAggregationFields on DepartmentAggregation {
		  slug
		  name
		  __typename
		}

	`
	expectedTreeJ := map[string][]string{
		"": {"search"},
		"search": {
			"aggregations", "some_field"},
		"search.aggregations": {"departments"},
		"search.aggregations.departments": {
			"slug", "name", "__typename"}}

	generatedTreeJ := BuildTreeUsingAliases(graphqlQueryJ, Variables{})

	assert.Equal(t, expectedTreeJ, generatedTreeJ)

	var graphqlQueryK string = `
		query ProductsSearchPage($include_aggregations: Boolean!) {
		  search {
		    aggregations @include(if:false) {
		      departments {
		        ...departmentAggregationFields
		        __typename
		      }
		    }
		    some_field
		  }
		}

		fragment departmentAggregationFields on DepartmentAggregation {
		  slug
		  name
		  __typename
		}

	`
	expectedTreeK := map[string][]string{
		"": {"search"},
		"search": {
			"aggregations_FALSE", "some_field"},
		"search.aggregations_FALSE": {"departments"},
		"search.aggregations_FALSE.departments": {
			"slug", "name", "__typename"}}

	generatedTreeK := BuildTreeUsingAliases(graphqlQueryK, Variables{})

	assert.Equal(t, expectedTreeK, generatedTreeK)

	var graphqlQueryL string = `
		query ProductsSearchPage($include_aggregations: Boolean!) {
		  search {
		    aggregations @include(if: $include_aggregations) {
		      departments {
		        ...departmentAggregationFields
		        __typename
		      }
		    }
		    some_field
		  }
		}

		fragment departmentAggregationFields on DepartmentAggregation {
		  slug
		  name
		  __typename
		}

	`
	expectedTreeL := map[string][]string{
		"": {"search"},
		"search": {
			"aggregations_FALSE", "some_field"},
		"search.aggregations_FALSE": {"departments"},
		"search.aggregations_FALSE.departments": {
			"slug", "name", "__typename"}}

	generatedTreeL := BuildTreeUsingAliases(graphqlQueryL, Variables{
		"include_aggregations": false})

	assert.Equal(t, expectedTreeL, generatedTreeL)

	var graphqlQueryM string = `
		query ProductsSearchPage($include_aggregations: Boolean!) {
		  search {
		    aggregations @include(if: $include_aggregations) {
		      departments {
		        ...departmentAggregationFields
		        __typename
		      }
		    }
		    lorem @include(if: $blabla) {

		    }
		    some_field
		  }
		}

		fragment departmentAggregationFields on DepartmentAggregation {
		  slug
		  name
		  __typename
		}

	`
	expectedTreeM := map[string][]string{
		"":                    {"search"},
		"search":              {"aggregations", "lorem", "some_field"},
		"search.aggregations": {"departments"},
		"search.aggregations.departments": {
			"slug", "name", "__typename"}}

	generatedTreeM := BuildTreeUsingAliases(graphqlQueryM, Variables{
		"include_aggregations": true,
		"blabla":               "lorem"})

	assert.Equal(t, expectedTreeM, generatedTreeM)
}
