package ezutil_test

import (
	"testing"
	"time"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model for GORM scopes testing
type TestModel struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:100"`
	Age       int       `gorm:"default:0"`
	Active    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	return db
}

func seedTestData(t *testing.T, db *gorm.DB) {
	testData := []TestModel{
		{Name: "Alice", Age: 25, Active: true, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "Bob", Age: 30, Active: false, CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Name: "Charlie", Age: 35, Active: true, CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
		{Name: "Diana", Age: 28, Active: false, CreatedAt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)},
		{Name: "Eve", Age: 32, Active: true, CreatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)},
	}

	for _, data := range testData {
		if err := db.Create(&data).Error; err != nil {
			t.Fatalf("Failed to seed test data: %v", err)
		}
	}
}

func TestPaginate(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	tests := []struct {
		name           string
		page           int
		limit          int
		expectedCount  int
		expectedOffset int
	}{
		{
			name:           "first page",
			page:           1,
			limit:          2,
			expectedCount:  2,
			expectedOffset: 0,
		},
		{
			name:           "second page",
			page:           2,
			limit:          2,
			expectedCount:  2,
			expectedOffset: 2,
		},
		{
			name:           "third page",
			page:           3,
			limit:          2,
			expectedCount:  1,
			expectedOffset: 4,
		},
		{
			name:           "page beyond data",
			page:           10,
			limit:          2,
			expectedCount:  0,
			expectedOffset: 18,
		},
		{
			name:           "zero page defaults to 1",
			page:           0,
			limit:          3,
			expectedCount:  3,
			expectedOffset: 0,
		},
		{
			name:           "negative page defaults to 1",
			page:           -1,
			limit:          3,
			expectedCount:  3,
			expectedOffset: 0,
		},
		{
			name:           "large limit",
			page:           1,
			limit:          10,
			expectedCount:  5,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			err := db.Scopes(ezutil.Paginate(tt.page, tt.limit)).Find(&results).Error
			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
		})
	}
}

func TestOrderBy(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	tests := []struct {
		name      string
		field     string
		ascending bool
		expectErr bool
		expected  []string // expected names in order
	}{
		{
			name:      "order by name ascending",
			field:     "name",
			ascending: true,
			expectErr: false,
			expected:  []string{"Alice", "Bob", "Charlie", "Diana", "Eve"},
		},
		{
			name:      "order by name descending",
			field:     "name",
			ascending: false,
			expectErr: false,
			expected:  []string{"Eve", "Diana", "Charlie", "Bob", "Alice"},
		},
		{
			name:      "order by age ascending",
			field:     "age",
			ascending: true,
			expectErr: false,
			expected:  []string{"Alice", "Diana", "Bob", "Eve", "Charlie"},
		},
		{
			name:      "order by age descending",
			field:     "age",
			ascending: false,
			expectErr: false,
			expected:  []string{"Charlie", "Eve", "Bob", "Diana", "Alice"},
		},
		{
			name:      "invalid field name",
			field:     "invalid; DROP TABLE users;",
			ascending: true,
			expectErr: true,
		},
		{
			name:      "field with special characters",
			field:     "field@name",
			ascending: true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			query := db.Scopes(ezutil.OrderBy(tt.field, tt.ascending))
			err := query.Find(&results).Error
			
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, 5)
				
				// Check order
				for i, expected := range tt.expected {
					assert.Equal(t, expected, results[i].Name)
				}
			}
		})
	}
}

func TestWhereBySpec(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	tests := []struct {
		name          string
		spec          TestModel
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "filter by active true",
			spec:          TestModel{Active: true},
			expectedCount: 3,
			expectedNames: []string{"Alice", "Charlie", "Eve"},
		},
		{
			name:          "filter by active false",
			spec:          TestModel{Active: false},
			expectedCount: 5, // GORM ignores false (zero value), so returns all
			expectedNames: []string{}, // Don't check names since all are returned
		},
		{
			name:          "filter by name",
			spec:          TestModel{Name: "Alice"},
			expectedCount: 1,
			expectedNames: []string{"Alice"},
		},
		{
			name:          "filter by age",
			spec:          TestModel{Age: 30},
			expectedCount: 1,
			expectedNames: []string{"Bob"},
		},
		{
			name:          "filter by multiple fields",
			spec:          TestModel{Age: 25, Active: true},
			expectedCount: 1,
			expectedNames: []string{"Alice"},
		},
		{
			name:          "no matching records",
			spec:          TestModel{Name: "NonExistent"},
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name:          "empty spec returns all",
			spec:          TestModel{},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			err := db.Scopes(ezutil.WhereBySpec(tt.spec)).Find(&results).Error
			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			
			if len(tt.expectedNames) > 0 {
				actualNames := make([]string, len(results))
				for i, result := range results {
					actualNames[i] = result.Name
				}
				assert.ElementsMatch(t, tt.expectedNames, actualNames)
			}
		})
	}
}

func TestPreloadRelations(t *testing.T) {
	// For this test, we'll create a more complex model with relations
	type Post struct {
		ID     uint   `gorm:"primarykey"`
		Title  string `gorm:"size:200"`
		UserID uint
	}

	type Profile struct {
		ID     uint   `gorm:"primarykey"`
		Bio    string `gorm:"size:500"`
		UserID uint
	}

	type User struct {
		ID       uint     `gorm:"primarykey"`
		Name     string   `gorm:"size:100"`
		Posts    []Post   `gorm:"foreignKey:UserID"`
		Profile  Profile  `gorm:"foreignKey:UserID"`
	}

	db := setupTestDB(t)
	err := db.AutoMigrate(&User{}, &Post{}, &Profile{})
	assert.NoError(t, err)

	// Seed data
	user := User{Name: "John"}
	db.Create(&user)
	
	post1 := Post{Title: "Post 1", UserID: user.ID}
	post2 := Post{Title: "Post 2", UserID: user.ID}
	db.Create(&post1)
	db.Create(&post2)
	
	profile := Profile{Bio: "John's bio", UserID: user.ID}
	db.Create(&profile)

	tests := []struct {
		name      string
		relations []string
	}{
		{
			name:      "preload single relation",
			relations: []string{"Posts"},
		},
		{
			name:      "preload multiple relations",
			relations: []string{"Posts", "Profile"},
		},
		{
			name:      "preload no relations",
			relations: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result User
			err := db.Scopes(ezutil.PreloadRelations(tt.relations)).First(&result).Error
			assert.NoError(t, err)
			assert.Equal(t, "John", result.Name)
			
			// Check if relations are loaded based on what was requested
			for _, relation := range tt.relations {
				switch relation {
				case "Posts":
					assert.Len(t, result.Posts, 2)
				case "Profile":
					assert.Equal(t, "John's bio", result.Profile.Bio)
				}
			}
		})
	}
}

func TestBetweenTime(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	startTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 4, 23, 59, 59, 0, time.UTC)
	zeroTime := time.Time{}

	tests := []struct {
		name          string
		col           string
		start         time.Time
		end           time.Time
		expectedCount int
	}{
		{
			name:          "between two dates",
			col:           "created_at",
			start:         startTime,
			end:           endTime,
			expectedCount: 3, // Bob, Charlie, Diana
		},
		{
			name:          "from start date only",
			col:           "created_at",
			start:         startTime,
			end:           zeroTime,
			expectedCount: 4, // Bob, Charlie, Diana, Eve
		},
		{
			name:          "until end date only",
			col:           "created_at",
			start:         zeroTime,
			end:           endTime,
			expectedCount: 4, // Alice, Bob, Charlie, Diana
		},
		{
			name:          "both dates zero",
			col:           "created_at",
			start:         zeroTime,
			end:           zeroTime,
			expectedCount: 5, // All records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			err := db.Scopes(ezutil.BetweenTime(tt.col, tt.start, tt.end)).Find(&results).Error
			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
		})
	}
}

func TestDefaultOrder(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	var results []TestModel
	err := db.Scopes(ezutil.DefaultOrder()).Find(&results).Error
	assert.NoError(t, err)
	assert.Len(t, results, 5)

	// Should be ordered by created_at DESC
	expectedOrder := []string{"Eve", "Diana", "Charlie", "Bob", "Alice"}
	for i, expected := range expectedOrder {
		assert.Equal(t, expected, results[i].Name)
	}
}

func TestForUpdate(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	tests := []struct {
		name   string
		enable bool
	}{
		{
			name:   "for update enabled",
			enable: true,
		},
		{
			name:   "for update disabled",
			enable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []TestModel
			err := db.Scopes(ezutil.ForUpdate(tt.enable)).Find(&results).Error
			assert.NoError(t, err)
			assert.Len(t, results, 5)
		})
	}
}

func TestCombinedScopes(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	// Test combining multiple scopes
	t.Run("paginate + order + filter", func(t *testing.T) {
		spec := TestModel{Active: true}
		var results []TestModel
		
		err := db.Scopes(
			ezutil.WhereBySpec(spec),
			ezutil.OrderBy("age", true), // ascending by age
			ezutil.Paginate(1, 2),       // first 2 results
		).Find(&results).Error
		
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		
		// Should get Alice (25) and Eve (32) - first 2 active users by age
		assert.Equal(t, "Alice", results[0].Name)
		assert.Equal(t, "Eve", results[1].Name)
	})

	t.Run("time range + order + pagination", func(t *testing.T) {
		startTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 4, 23, 59, 59, 0, time.UTC)
		
		var results []TestModel
		err := db.Scopes(
			ezutil.BetweenTime("created_at", startTime, endTime),
			ezutil.OrderBy("name", false), // descending by name
			ezutil.Paginate(1, 2),          // first 2 results
		).Find(&results).Error
		
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		
		if len(results) >= 2 {
			// Should get Diana and Charlie (first 2 in desc name order within date range)
			assert.Equal(t, "Diana", results[0].Name)
			assert.Equal(t, "Charlie", results[1].Name)
		}
	})
}
