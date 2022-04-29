package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_CreateTag_succeeds_addsTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)

	resolver := ResolverForTag{DB: db.DB}
	tag, err := resolver.CreateTag(fake.User(5), "new tag", "#fff")

	require.Nil(t, err)
	expected := &gqlmodel.TagDefinition{
		Key:   "new tag",
		Color: "#fff",
	}
	require.Equal(t, expected, tag)
	assertTagExist(t, db, model.TagDefinition{
		Key:    "new tag",
		Color:  "#fff",
		UserID: 5,
	})
	assertTagCount(t, db, 1)
}

func TestGQL_CreateTag_fails_tagAlreadyExists(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TagDefinition{Key: "existing tag", Color: "#fff", UserID: 5})

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.CreateTag(fake.User(5), "existing tag", "#fff")

	require.EqualError(t, err, "tag with key 'existing tag' does already exist")
	assertTagCount(t, db, 1)
}

func TestGQL_CreateTag_fails_tagAlreadyExists_caseInsensitive(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TagDefinition{Key: "tag", Color: "#fff", UserID: 5})

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.CreateTag(fake.User(5), "Tag", "#fff")

	require.EqualError(t, err, "tag with key 'tag' does already exist")
	assertTagCount(t, db, 1)
}

func TestGQL_CreateTag_succeeds_existingTagForOtherUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(4)
	db.User(5)
	db.Create(&model.TagDefinition{Key: "existing tag", Color: "#fff", UserID: 4})

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.CreateTag(fake.User(5), "existing tag", "#xxx")

	assert.Nil(t, err)
	assertTagCount(t, db, 2)
	assertTagExist(t, db, model.TagDefinition{
		Key:    "existing tag",
		Color:  "#xxx",
		UserID: 5,
	})
}
