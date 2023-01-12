package myra_search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/raralabs/myra-search/pkg/models"
	"github.com/raralabs/myra-search/pkg/utils"
	"github.com/raralabs/myra-search/pkg/utils/db/pgdb"
	"gorm.io/gorm"
)

// ClientInterface exposes the needed methods to external usage
type ClientInterface interface {
	Search(slug, search string, pagination ...int) ([]models.ResponseSearchIndex, error)
	SearchByField(slug string, fieldSearch map[string]interface{}, pagination ...int) ([]models.ResponseSearchIndex, error)
	Index(slug string, uid string, tableInfo string, action map[string]interface{}, searchValue map[string]interface{}) (string, error)
	Delete(slug, uid, tableInfo string) (string, error)
	CloseConnection()
}

type Client struct {
	db *gorm.DB
}

// NewClient initializes the Client and database and returns the ClientInterface instance.
func NewClient(dsn string) ClientInterface {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return Client{
		db: pgdb.ConnectDatabase(dsn),
	}
}

// Search takes the slug and search string as input for global searching the value and returns the multiple matching recoreds.
func (s Client) Search(slug, search string, pagination ...int) ([]models.ResponseSearchIndex, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return []models.ResponseSearchIndex{}, errors.New("DB Connection Failed")
	}
	offset, limit := utils.Pagination(pagination...)
	search = strings.TrimSpace(search)
	search = strings.Join(strings.Fields(search), " ")
	if search == "" {
		return []models.ResponseSearchIndex{}, errors.New("please provide the search string")
	}
	search = strings.ReplaceAll(search, " ", ":*&")
	var model []models.ResponseSearchIndex
	query := fmt.Sprintf("SELECT id, table_info, action_info FROM \"%s\".search_indices ", slug)
	err := s.db.Raw(query+" WHERE tsv_text @@ to_tsquery(? || ':*') ORDER BY id OFFSET ? LIMIT ?", search, offset, limit).Scan(&model).Error
	return model, err
}

// SearchByField search into the field level and return the needed action.
func (s Client) SearchByField(slug string, fieldSearch map[string]interface{}, pagination ...int) ([]models.ResponseSearchIndex, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return []models.ResponseSearchIndex{}, errors.New("DB Connection Failed")
	}
	offset, limit := utils.Pagination(pagination...)
	var model []models.ResponseSearchIndex
	query := fmt.Sprintf("SELECT id, table_info, action_info FROM \"%s\".search_indices WHERE ", slug)
	i := 0
	len := len(fieldSearch)
	for k, v := range fieldSearch {
		if len < i {
			query += fmt.Sprintf("search_field->>'%s' like '%v%s' and ", k, v, "%")
		} else {
			query += fmt.Sprintf("search_field->>'%s' like '%v%s' ", k, v, "%")
		}
	}
	err := s.db.Raw(query+" OFFSET ? LIMIT ?", offset, limit).Scan(&model).Error
	return model, err
}

// Index takes the slug, uid, table_info, action, search_value as input to create the index in the database.
func (s Client) Index(slug string, uid string, tableInfo string, action map[string]interface{}, searchValue map[string]interface{}) (string, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return "", errors.New("DB Connection Failed")
	}
	tsv := ""
	first := true
	for _, value := range searchValue {
		if first {
			tsv += fmt.Sprintf("%v", value)
			first = false
		} else {
			tsv += fmt.Sprintf(" %v", value)
		}
	}
	query := fmt.Sprintf("INSERT INTO \"%s\".search_indices(id,table_info,action_info,tsv_text, search_field)", slug)
	var id string
	err := s.db.
		Raw(query+" VALUES(?,?,?,to_tsvector(?),?) ON CONFLICT (id,table_info) DO UPDATE SET action_info=?, search_field=?, tsv_text=to_tsvector(?) RETURNING id",
			uid,
			tableInfo,
			utils.MapToJSON(action),
			tsv,
			utils.MapToJSON(searchValue),
			utils.MapToJSON(action),
			utils.MapToJSON(searchValue),
			tsv,
		).
		Scan(&id).Error
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s Client) Delete(slug, uid, tableInfo string) (string, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return "", errors.New("DB Connection Failed")
	}
	var id string
	query := fmt.Sprintf("DELETE FROM \"%s\".search_indices", slug)
	err := s.db.Raw(query+" WHERE id = ? and table_info = ? RETURNING id", uid, tableInfo).Scan(&id).Error
	if err != nil || id == "" {
		return id, errors.New("record not found")
	}
	return id, nil
}

func (s Client) CloseConnection() {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return
	}
	pgdb.CloseConnection(s.db)
}
