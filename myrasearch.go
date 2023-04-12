package tsv_search_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/raralabs/tsv-search-interface/pkg/models"
	"github.com/raralabs/tsv-search-interface/pkg/utils"
	"github.com/raralabs/tsv-search-interface/pkg/utils/db/pgdb"
	"gorm.io/gorm"
	"math"
	"reflect"
	"strings"
)

// ClientInterface exposes the needed methods to external usage
type ClientInterface interface {
	Search(slug, search string, pagination ...int) ([]models.ResponseSearchIndex, error)
	InternalSearch(slug, search string, tableInfo string, pagination ...int) ([]string, error)
	SearchByField(slug string, fieldSearch map[string]interface{}, pagination ...int) ([]models.ResponseSearchIndex, error)
	Index(slug string, uid string, tableInfo string, action map[string]interface{}, searchValue map[string]interface{}) (string, error)
	IndexInternal(slug string, uid string, tableInfo string, searchValue map[string]interface{}) (string, error)
	IndexBatchInternal(slug string, tableInfo string, input []models.BatchIndexInput) error
	Delete(slug, uid, tableInfo string) (string, error)
	CloseConnection()
}

type Client struct {
	db *gorm.DB
}

const BatchSize = 10000

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

// InternalSearch takes the slug and search string as input for global searching the value and returns the multiple matching recoreds.
func (s Client) InternalSearch(slug, search string, tableInfo string, pagination ...int) ([]string, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return []string{}, errors.New("DB Connection Failed")
	}
	offset, limit := utils.Pagination(pagination...)
	search = strings.TrimSpace(search)
	search = strings.Join(strings.Fields(search), " ")
	if search == "" {
		return []string{}, errors.New("please provide the search string")
	}
	search = strings.ReplaceAll(search, " ", ":*&")
	var model []string
	query := fmt.Sprintf("SELECT id FROM \"%s\".internal_search_indices ", slug)
	err := s.db.Raw(query+" WHERE table_info = ? and tsv_text @@ to_tsquery(? || ':*') ORDER BY id OFFSET ? LIMIT ?", tableInfo, search, offset, limit).Scan(&model).Error
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

func getTableList(s Client, counter int, tableInfo ...interface{}) []models.RelatedInfoWithLevel {
	var rF []models.RelatedInfoWithLevel
	if counter >= 5 || len(tableInfo) == 0 {
		return rF
	}
	l := gl(s, tableInfo...)
	rF = append(rF, models.RelatedInfoWithLevel{RelatedInfo: l, Counter: counter})
	var tF []interface{}
	for _, v := range l {
		tF = append(tF, v.RelatedTable)
	}
	if len(tF) == 0 {
		return rF
	}
	return append(rF, getTableList(s, counter+1, tF...)...)
}

func gl(s Client, tableInfo ...interface{}) []models.RelatedInfo {
	var tableList []models.RelatedInfo
	query := "SELECT * FROM related_infos WHERE table_info = ?"
	for i := 1; i < len(tableInfo); i++ {
		query += " or table_info = ? "
	}
	s.db.Raw(query, tableInfo...).Scan(&tableList)
	return tableList
}

func getTableInfo(s Client, tableInfo string) models.TableInformation {
	var tableInformation models.TableInformation
	s.db.Raw("SELECT * FROM table_information WHERE table_name = ? limit 1", tableInfo).Scan(&tableInformation)
	return tableInformation
}

func skip(value interface{}, skipId bool) bool {
	switch value {
	case "created_at", "modified_at", "creator_id", "modifier_id", "deleted_at":
		return true
	case "id":
		return skipId
	default:
		return false
	}
}

// IndexInternal takes the slug, uid, table_info, search_value as input to create the index in the database.
func (s Client) IndexInternal(slug string, uid string, tableInfo string, searchValue map[string]interface{}) (string, error) {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return "", errors.New("DB Connection Failed")
	}
	tableInformation := getTableInfo(s, tableInfo)
	tsv := ""
	first := true
	tempSearchValue := map[string]interface{}{}
	for key, value := range searchValue {
		if reflect.TypeOf(value).Kind() == reflect.String {
			value = strings.ReplaceAll(value.(string), "\u0000", "")
			searchValue[key] = value
		}
		if value == "" || skip(key, false) || tableInformation.TableName == "" || !strings.Contains(tableInformation.ColumnName, fmt.Sprintf("%v", key)) {
			continue
		}
		if first {
			tsv += fmt.Sprintf("%v", value)
			first = false
		} else {
			tsv += fmt.Sprintf(" %v", value)
		}
		tempSearchValue[fmt.Sprintf("%s.%s", tableInfo, key)] = value
	}

	tL := getTableList(s, 0, tableInfo)
	if len(tL) > 0 {
		for _, t := range tL {
			for _, value1 := range t.RelatedInfo {
				if term, ok := tempSearchValue[fmt.Sprintf("%s.%s", value1.TableInfo, value1.ForeignField)]; ok {
					var internalSearch models.InternalSearchIndex
					query := fmt.Sprintf("SELECT search_field FROM \"%s\".internal_search_indices ", slug)
					if strings.ToLower(value1.MappingField) != "id" {
						s.db.Raw(query+" WHERE table_info=? and search_field ->> ?  = ? order by id desc limit 1", value1.RelatedTable, value1.MappingField, term).Scan(&internalSearch)
					} else {
						s.db.Raw(query+" WHERE table_info=? and id = ?", value1.RelatedTable, term).Scan(&internalSearch)
					}
					tableInformation := getTableInfo(s, value1.RelatedTable)
					d, err := internalSearch.SearchField.MarshalJSON()
					if err != nil {
						continue
					}
					data := map[string]interface{}{}
					err = json.Unmarshal(d, &data)
					if err != nil {
						continue
					}
					for key, value := range data {
						if value == "" || skip(key, false) || tableInformation.TableName == "" || !strings.Contains(tableInformation.ColumnName, fmt.Sprintf("%v", key)) {
							continue
						}
						if first {
							tsv += fmt.Sprintf("%v", value)
							first = false
						} else {
							tsv += fmt.Sprintf(" %v", value)
						}
						tempSearchValue[fmt.Sprintf("%s.%s", tableInformation.TableName, key)] = value
					}
				}
			}
		}
	}

	if tsv == "" {
		return "", nil
	}
	query := fmt.Sprintf("INSERT INTO \"%s\".internal_search_indices(id,table_info,tsv_text, search_field)", slug)
	var id string
	err := s.db.
		Raw(query+" VALUES(?,?,to_tsvector(?),?) ON CONFLICT (id,table_info) DO UPDATE SET search_field=?, tsv_text=to_tsvector(?) RETURNING id",
			uid,
			tableInfo,
			tsv,
			utils.MapToJSON(searchValue),
			utils.MapToJSON(searchValue),
			tsv,
		).
		Scan(&id).Error
	if err != nil {
		return "", err
	}
	return id, nil
}

// IndexBatchInternal takes the slug, uid, table_info, search_value as input to create the index in the database.
func (s Client) IndexBatchInternal(slug string, tableInfo string, input []models.BatchIndexInput) error {
	if s.db == nil {
		fmt.Errorf("%v", errors.New("DB Connection Failed"))
		return errors.New("DB Connection Failed")
	}
	tableInformation := getTableInfo(s, tableInfo)
	page := int64(math.Ceil(float64(len(input)) / float64(BatchSize)))
	for i := int64(0); i < page; i++ {
		var slice []models.BatchIndexInput
		if i == page-1 {
			slice = input[i*BatchSize:]
		} else {
			slice = input[i*BatchSize : (i+1)*BatchSize]
		}
		valueStrings := make([]string, 0, BatchSize)
		valueArgs := make([]interface{}, 0, BatchSize*4)
		for _, sv := range slice {
			tsv := ""
			first := true
			tempSearchValue := map[string]interface{}{}
			for key, value := range sv.SearchValue {
				if reflect.TypeOf(value).Kind() == reflect.String {
					value = strings.ReplaceAll(value.(string), "\u0000", "")
					sv.SearchValue[key] = value
				}
				if value == "" || skip(key, false) || tableInformation.TableName == "" || !strings.Contains(tableInformation.ColumnName, fmt.Sprintf("%v", key)) {
					continue
				}
				if first {
					tsv += fmt.Sprintf("%v", value)
					first = false
				} else {
					tsv += fmt.Sprintf(" %v", value)
				}
				tempSearchValue[fmt.Sprintf("%s.%s", tableInfo, key)] = value
			}

			tL := getTableList(s, 0, tableInfo)
			if len(tL) > 0 {
				for _, t := range tL {
					for _, value1 := range t.RelatedInfo {
						if term, ok := tempSearchValue[fmt.Sprintf("%s.%s", value1.TableInfo, value1.ForeignField)]; ok {
							var internalSearch models.InternalSearchIndex
							query := fmt.Sprintf("SELECT search_field FROM \"%s\".internal_search_indices ", slug)
							if strings.ToLower(value1.MappingField) != "id" {
								s.db.Raw(query+" WHERE table_info=? and search_field ->> ?  = ? order by id desc limit 1", value1.RelatedTable, value1.MappingField, term).Scan(&internalSearch)
							} else {
								s.db.Raw(query+" WHERE table_info=? and id = ?", value1.RelatedTable, term).Scan(&internalSearch)
							}
							tableInformation := getTableInfo(s, value1.RelatedTable)
							d, err := internalSearch.SearchField.MarshalJSON()
							if err != nil {
								continue
							}
							data := map[string]interface{}{}
							err = json.Unmarshal(d, &data)
							if err != nil {
								continue
							}
							for key, value := range data {
								if value == "" || skip(key, false) || tableInformation.TableName == "" || !strings.Contains(tableInformation.ColumnName, fmt.Sprintf("%v", key)) {
									continue
								}
								if first {
									tsv += fmt.Sprintf("%v", value)
									first = false
								} else {
									tsv += fmt.Sprintf(" %v", value)
								}
								tempSearchValue[fmt.Sprintf("%s.%s", tableInformation.TableName, key)] = value
							}
						}
					}
				}
			}

			if tsv == "" {
				continue
			}
			valueStrings = append(valueStrings, "(?,?,to_tsvector(?),?)")
			valueArgs = append(valueArgs, sv.UID)
			valueArgs = append(valueArgs, tableInfo)
			valueArgs = append(valueArgs, tsv)
			valueArgs = append(valueArgs, utils.MapToJSON(sv.SearchValue))
		}
		if len(valueStrings) == 0 {
			continue
		}
		query := fmt.Sprintf("INSERT INTO \"%s\".internal_search_indices(id,table_info,tsv_text, search_field) VALUES %s ", slug, strings.Join(valueStrings, ","))
		var id string
		err := s.db.
			Raw(query+" ON CONFLICT (id,table_info) DO UPDATE SET tsv_text = EXCLUDED.tsv_text, search_field = EXCLUDED.search_field", valueArgs...).
			Scan(&id).Error
		if err != nil {
			return err
		}
	}
	return nil
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
