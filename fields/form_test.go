package fields

import (
	. "testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func TestJSONStatement(t *T) {
	f := &Form{
		Fields: []interface{}{
			&Statement{
				Field: Field{
					Type: TypeStatement,
				},
			},
		},
	}
	fs := `{"title":"","fields":[{"type":"statement","question":""}]}`
	j, err := json.Marshal(f)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	nf := &Form{}
	err = json.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}

func TestBSONStatement(t *T) {
	s := &Statement{
		Field: Field{
			Type: TypeStatement,
			Question: "Hey?",
		},
	}
	f := &Form{
		Fields: []interface{}{s},
	}
	fexp := struct{
		Title string `bson:"t"`
		Fields []*Statement `bson:"f"`
	}{
		Title: "",
		Fields: []*Statement{s},
	}
	j, err := bson.Marshal(f)
	require.Nil(t, err)
	jexp, err := bson.Marshal(fexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	nf := &Form{}
	err = bson.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}

func TestJSONMultipleChoice(t *T) {
	f := &Form{
		Fields: []interface{}{
			&MultipleChoice{
				Field: Field{
					Type: TypeMultipleChoice,
				},
				Choices: []MultipleChoiceChoice{
					MultipleChoiceChoice{
						Label: "Label",
					},
				},
			},
		},
	}
	fs := `{"title":"","fields":[{"type":"multiple_choice","question":"","choices":[{"label":"Label"}]}]}`
	j, err := json.Marshal(f)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	nf := &Form{}
	err = json.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}

func TestBSONMultipleChoice(t *T) {
	mc := &MultipleChoice{
		Field: Field{
			Type: TypeMultipleChoice,
		},
		Choices: []MultipleChoiceChoice{
			MultipleChoiceChoice{
				Label: "Label",
			},
		},
	}
	f := &Form{
		Fields: []interface{}{mc},
	}
	fexp := struct {
		Title  string `bson:"t"`
		Fields []*MultipleChoice `bson:"f"`
	}{
		Title: "",
		Fields: []*MultipleChoice{mc},
	}
	j, err := bson.Marshal(f)
	require.Nil(t, err)
	jexp, err := bson.Marshal(fexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	nf := &Form{}
	err = bson.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}

func TestJSONOpinionLabels(t *T) {
	f := &Form{
		Fields: []interface{}{
			&OpinionScale{
				Field: Field{
					Type: TypeOpinionScale,
				},
				Steps: 5,
				StartAtOne: true,
				Labels: OpinionLabels{
					Left: "l",
					Center: "c",
					Right: "r",
				},
			},
		},
	}
	fs := `{"title":"","fields":[{"type":"opinion_scale","question":"","steps":5,"start_at_one":true,"labels":{"left":"l","center":"c","right":"r"}}]}`
	j, err := json.Marshal(f)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	nf := &Form{}
	err = json.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}

func TestBSONOpinionLabels(t *T) {
	os := &OpinionScale{
		Field: Field{
			Type: TypeOpinionScale,
		},
		Steps: 5,
		StartAtOne: true,
		Labels: OpinionLabels{
			Left: "l",
			Center: "c",
			Right: "r",
		},
	}
	f := &Form{
		Fields: []interface{}{os},
	}
	fexp := struct {
		Title  string `bson:"t"`
		Fields []*OpinionScale `bson:"f"`
	}{
		Title: "",
		Fields: []*OpinionScale{os},
	}
	j, err := bson.Marshal(f)
	require.Nil(t, err)
	jexp, err := bson.Marshal(fexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	nf := &Form{}
	err = bson.Unmarshal(j, nf)
	require.Nil(t, err)
	assert.EqualValues(t, f, nf)
}
