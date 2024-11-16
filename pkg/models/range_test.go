package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestRangeSurrealString(t *testing.T) {
	tests := map[string]*models.Range[int]{
		"1..=3": models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundIncluded(3)),
		"1>..3": models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundExcluded(3)),

		"1..3": models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundExcluded(3)),
		"1..":  models.NewRange[int](models.NewBoundIncluded(1), nil),
		"..3":  models.NewRange[int](nil, models.NewBoundExcluded(3)),

		"1>..=3": models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundIncluded(3)),
		"1>..":   models.NewRange[int](models.NewBoundExcluded(1), nil),
		"..=3":   models.NewRange[int](nil, models.NewBoundIncluded(3)),

		"..": models.NewRange[int](nil, nil),
	}
	for expected, src := range tests {
		if s, err := src.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, s)
		}
	}
}

func TestRangeCBOR(t *testing.T) {
	tests := []struct {
		src *models.Range[int]
		beg func(b *models.Bound[int])
		end func(b *models.Bound[int])
	}{
		// "1..=3"
		{
			src: models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundIncluded(3)),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
		},
		// "1>..3"
		{
			src: models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundExcluded(3)),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
		},

		// "1..3"
		{
			src: models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundExcluded(3)),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
		},
		// "1.."
		{
			src: models.NewRange[int](models.NewBoundIncluded(1), nil),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
		},
		// "..3"
		{
			src: models.NewRange[int](nil, models.NewBoundExcluded(3)),
			beg: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
		},

		// "1>..=3"
		{
			src: models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundIncluded(3)),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
		},
		// "1>.."
		{
			src: models.NewRange[int](models.NewBoundExcluded(1), nil),
			beg: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 1)
					assert.True(t, b.Excluded())
					assert.False(t, b.Included())
				}
			},
			end: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
		},
		// "..=3"
		{
			src: models.NewRange[int](nil, models.NewBoundIncluded(3)),
			beg: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
			end: func(b *models.Bound[int]) {
				if assert.NotNil(t, b) {
					assert.Equal(t, b.Value, 3)
					assert.False(t, b.Excluded())
					assert.True(t, b.Included())
				}
			},
		},

		// ".."
		{
			src: models.NewRange[int](nil, nil),
			beg: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
			end: func(b *models.Bound[int]) {
				assert.Nil(t, b)
			},
		},
	}
	for _, tt := range tests {
		data, err := models.CBORFormatter.Marshal(tt.src)
		if assert.NoError(t, err) {
			var dst models.Range[int]
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				tt.beg(dst.Begin)
				tt.end(dst.End)
			}
		}
	}
}

func TestRangeJSON(t *testing.T) {
	tests := map[string]*models.Range[int]{
		"1..=3": models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundIncluded(3)),
		"1>..3": models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundExcluded(3)),

		"1..3": models.NewRange[int](models.NewBoundIncluded(1), models.NewBoundExcluded(3)),
		"1..":  models.NewRange[int](models.NewBoundIncluded(1), nil),
		"..3":  models.NewRange[int](nil, models.NewBoundExcluded(3)),

		"1>..=3": models.NewRange[int](models.NewBoundExcluded(1), models.NewBoundIncluded(3)),
		"1>..":   models.NewRange[int](models.NewBoundExcluded(1), nil),
		"..=3":   models.NewRange[int](nil, models.NewBoundIncluded(3)),

		"..": models.NewRange[int](nil, nil),
	}
	for expected, src := range tests {
		data, err := models.JSONFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst string
			if err := models.JSONFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				assert.Equal(t, expected, dst)
			}
		}
	}
}
