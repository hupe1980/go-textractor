package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterRelationshipIDsByType(t *testing.T) {
	// Create a sample block with relationships
	block := types.Block{
		Relationships: []types.Relationship{
			{Type: "friend", Ids: []string{"id1", "id2"}},
			{Type: "family", Ids: []string{"id3", "id4"}},
			{Type: "friend", Ids: []string{"id5", "id6"}},
			{Type: "colleague", Ids: []string{"id7", "id8"}},
		},
	}

	// Test case 1: Filter for "friend" relationships
	t.Run("FilterFriendRelationships", func(t *testing.T) {
		expected := []string{"id1", "id2", "id5", "id6"}
		result := filterRelationshipIDsByType(block, "friend")
		assert.ElementsMatch(t, expected, result, "Failed to filter friend relationships")
	})

	// Test case 2: Filter for "family" relationships
	t.Run("FilterFamilyRelationships", func(t *testing.T) {
		expected := []string{"id3", "id4"}
		result := filterRelationshipIDsByType(block, "family")
		assert.ElementsMatch(t, expected, result, "Failed to filter family relationships")
	})

	// Test case 3: Filter for "colleague" relationships
	t.Run("FilterColleagueRelationships", func(t *testing.T) {
		expected := []string{"id7", "id8"}
		result := filterRelationshipIDsByType(block, "colleague")
		assert.ElementsMatch(t, expected, result, "Failed to filter colleague relationships")
	})

	// Test case 4: Filter for non-existent relationship type
	t.Run("FilterNonExistentRelationships", func(t *testing.T) {
		expected := []string{}
		result := filterRelationshipIDsByType(block, "nonexistent")
		assert.ElementsMatch(t, expected, result, "Failed to handle non-existent relationship type")
	})
}
