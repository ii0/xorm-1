package xorm

import (
	"database/sql"
	"strings"
	"time"
)

type SqlsExecutor struct {
	session *Session
	sqls    interface{}
	parmas  interface{}
	err     error
}

func (sqlsExecutor *SqlsExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	if sqlsExecutor.err != nil {
		return nil, nil, sqlsExecutor.err
	}
	var model_1_results ResultMap
	var model_2_results sql.Result
	var err error

	sqlModel := 1

	switch sqlsExecutor.sqls.(type) {
	case string:
		sqlStr := strings.TrimSpace(sqlsExecutor.sqls.(string))
		sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

		if sqlsExecutor.parmas == nil {
			switch sqlCmd {
			case "select", "desc":
				model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()

			case "insert", "delete", "update", "create":
				model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
				sqlModel = 2
			}
		} else {
			switch sqlsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, ok := sqlsExecutor.parmas.([]map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.Engine.AddSql(key, sqlStr)
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap[0]).Query()

				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap[0]).Execute()
					sqlModel = 2

				}
				sqlsExecutor.session.Engine.RemoveSql(key)
			case map[string]interface{}:
				parmaMap, ok := sqlsExecutor.parmas.(map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.Engine.AddSql(key, sqlStr)
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Execute()
					sqlModel = 2
				}
				sqlsExecutor.session.Engine.RemoveSql(key)
			default:
				return nil, nil, ErrParamsType
			}
		}

		resultSlice := make([][]map[string]interface{}, 1)

		if sqlModel == 1 {
			if model_1_results.Error != nil {
				return nil, nil, model_1_results.Error
			}

			resultSlice[0] = make([]map[string]interface{}, len(model_1_results.Results))
			resultSlice[0] = model_1_results.Results
			return resultSlice, nil, nil
		} else {
			if err != nil {
				return nil, nil, err
			}

			resultMap := make([]map[string]interface{}, 1)
			resultMap[0] = make(map[string]interface{})

			//todo all database support LastInsertId
			LastInsertId, _ := model_2_results.LastInsertId()

			resultMap[0]["LastInsertId"] = LastInsertId
			RowsAffected, err := model_2_results.RowsAffected()
			if err != nil {
				return nil, nil, err
			}
			resultMap[0]["RowsAffected"] = RowsAffected
			resultSlice[0] = resultMap
			return resultSlice, nil, nil
		}
	case []string:
		if sqlsExecutor.session.IsSqlFuc == true {
			err := sqlsExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlsSlice := sqlsExecutor.sqls.([]string)
		n := len(sqlsSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)
		switch sqlsExecutor.parmas.(type) {
		case []map[string]interface{}:
			parmaSlice = sqlsExecutor.parmas.([]map[string]interface{})

		default:
			if sqlsExecutor.session.IsSqlFuc == true {
				err := sqlsExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for i, _ := range sqlsSlice {
			sqlStr := strings.TrimSpace(sqlsSlice[i])
			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmaSlice[i] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()

				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
					sqlModel = 2
				}
			} else {
				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.Engine.AddSql(key, sqlStr)
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaSlice[i]).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaSlice[i]).Execute()
					sqlModel = 2
				}
				sqlsExecutor.session.Engine.RemoveSql(key)
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlsExecutor.session.IsSqlFuc == true {
						err := sqlsExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, model_1_results.Error
				}

				resultSlice[i] = make([]map[string]interface{}, len(model_1_results.Results))
				resultSlice[i] = model_1_results.Results

			} else {
				if err != nil {
					if sqlsExecutor.session.IsSqlFuc == true {
						err := sqlsExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}

				resultMap := make([]map[string]interface{}, 1)
				resultMap[0] = make(map[string]interface{})

				//todo all database support LastInsertId
				LastInsertId, _ := model_2_results.LastInsertId()

				resultMap[0]["LastInsertId"] = LastInsertId
				RowsAffected, err := model_2_results.RowsAffected()
				if err != nil {
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultSlice[i] = make([]map[string]interface{}, 1)
				resultSlice[i] = resultMap

			}
		}

		if sqlsExecutor.session.IsSqlFuc == true {
			err := sqlsExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return resultSlice, nil, nil

	case map[string]string:
		if sqlsExecutor.session.IsSqlFuc == true {
			err := sqlsExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlsMap := sqlsExecutor.sqls.(map[string]string)
		n := len(sqlsMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)
		switch sqlsExecutor.parmas.(type) {
		case map[string]map[string]interface{}:
			parmasMap = sqlsExecutor.parmas.(map[string]map[string]interface{})

		default:
			if sqlsExecutor.session.IsSqlFuc == true {
				err := sqlsExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for k, _ := range sqlsMap {
			sqlStr := strings.TrimSpace(sqlsMap[k])
			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmasMap[k] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()

				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
					sqlModel = 2
				}
			} else {
				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.Engine.AddSql(key, sqlStr)
				parmaMap := parmasMap[k]
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Execute()
					sqlModel = 2
				}
				sqlsExecutor.session.Engine.RemoveSql(key)
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlsExecutor.session.IsSqlFuc == true {
						err := sqlsExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, model_1_results.Error
				}

				resultsMap[k] = make([]map[string]interface{}, len(model_1_results.Results))
				resultsMap[k] = model_1_results.Results

			} else {
				if err != nil {
					if sqlsExecutor.session.IsSqlFuc == true {
						err := sqlsExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}

				resultMap := make([]map[string]interface{}, 1)
				resultMap[0] = make(map[string]interface{})

				//todo all database support LastInsertId
				LastInsertId, _ := model_2_results.LastInsertId()

				resultMap[0]["LastInsertId"] = LastInsertId
				RowsAffected, err := model_2_results.RowsAffected()
				if err != nil {
					if sqlsExecutor.session.IsSqlFuc == true {
						err := sqlsExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultsMap[k] = make([]map[string]interface{}, 1)
				resultsMap[k] = resultMap

			}
		}
		if sqlsExecutor.session.IsSqlFuc == true {
			err := sqlsExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
