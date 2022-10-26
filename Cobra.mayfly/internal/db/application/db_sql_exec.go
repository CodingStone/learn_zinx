package application

import (
	"encoding/json"
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/db/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/db/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"

	"strings"

	"github.com/xwb1989/sqlparser"
)

type DbSqlExecReq struct {
	DbId         uint64
	Db           string
	Sql          string
	Remark       string
	LoginAccount *model.LoginAccount
	DbInstance   *DbInstance
}

type DbSqlExecRes struct {
	ColNames []string
	Res      []map[string]interface{}
}

type DbSqlExec interface {
	// 执行sql
	Exec(execSqlReq *DbSqlExecReq) (*DbSqlExecRes, error)

	// 根据条件删除sql执行记录
	DeleteBy(condition *entity.DbSqlExec)

	// 分页获取
	GetPageList(condition *entity.DbSqlExec, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult
}

func newDbSqlExecApp(dbExecSqlRepo repository.DbSqlExec) DbSqlExec {
	return &dbSqlExecAppImpl{
		dbSqlExecRepo: dbExecSqlRepo,
	}
}

type dbSqlExecAppImpl struct {
	dbSqlExecRepo repository.DbSqlExec
}

func (d *dbSqlExecAppImpl) Exec(execSqlReq *DbSqlExecReq) (*DbSqlExecRes, error) {
	sql := execSqlReq.Sql
	loginAccount := execSqlReq.LoginAccount

	// 如果是pgsql并且为查询语句，则直接返回，sqlparser不支持解析pgsql select语句
	if execSqlReq.DbInstance.Type == entity.DbTypePostgres && (strings.HasPrefix(sql, "SELECT") || strings.HasPrefix(sql, "select")) {
		return doSelect(nil, execSqlReq, nil)
	}

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		global.Log.Error("sql解析失败: ", err)
		return nil, err
	}

	dbSqlExecRecord := new(entity.DbSqlExec)
	dbSqlExecRecord.DbId = execSqlReq.DbId
	dbSqlExecRecord.Db = execSqlReq.Db
	dbSqlExecRecord.Sql = sql
	dbSqlExecRecord.Remark = execSqlReq.Remark
	dbSqlExecRecord.SetBaseInfo(loginAccount)

	var execRes *DbSqlExecRes
	isSelect := false
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		isSelect = true
		execRes, err = doSelect(stmt, execSqlReq, dbSqlExecRecord)
	case *sqlparser.Show:
		isSelect = true
		execRes, err = doSelect(nil, execSqlReq, dbSqlExecRecord)
	case *sqlparser.OtherRead:
		isSelect = true
		execRes, err = doSelect(nil, execSqlReq, dbSqlExecRecord)
	case *sqlparser.Update:
		execRes, err = doUpdate(stmt, execSqlReq, dbSqlExecRecord)
	case *sqlparser.Delete:
		execRes, err = doDelete(stmt, execSqlReq, dbSqlExecRecord)
	case *sqlparser.Insert:
		execRes, err = doInsert(stmt, execSqlReq, dbSqlExecRecord)
	default:
		execRes, err = doExec(execSqlReq.Sql, execSqlReq.DbInstance)
	}
	if err != nil {
		return nil, err
	}

	if !isSelect {
		// 保存执行记录
		d.dbSqlExecRepo.Insert(dbSqlExecRecord)
	}
	return execRes, nil
}

func (d *dbSqlExecAppImpl) DeleteBy(condition *entity.DbSqlExec) {
	d.dbSqlExecRepo.DeleteBy(condition)
}

func (d *dbSqlExecAppImpl) GetPageList(condition *entity.DbSqlExec, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return d.dbSqlExecRepo.GetPageList(condition, pageParam, toEntity, orderBy...)
}

func doSelect(selectStmt *sqlparser.Select, execSqlReq *DbSqlExecReq, dbSqlExec *entity.DbSqlExec) (*DbSqlExecRes, error) {
	dbInstance := execSqlReq.DbInstance
	sql := execSqlReq.Sql
	colNames, res, err := dbInstance.SelectData(sql)
	if err != nil {
		return nil, err
	}
	return &DbSqlExecRes{
		ColNames: colNames,
		Res:      res,
	}, nil
}

func doUpdate(update *sqlparser.Update, execSqlReq *DbSqlExecReq, dbSqlExec *entity.DbSqlExec) (*DbSqlExecRes, error) {
	dbInstance := execSqlReq.DbInstance

	tableStr := sqlparser.String(update.TableExprs)
	// 可能使用别名，故空格切割
	tableName := strings.Split(tableStr, " ")[0]
	where := sqlparser.String(update.Where)

	updateExprs := update.Exprs
	updateColumns := make([]string, 0)
	for _, v := range updateExprs {
		updateColumns = append(updateColumns, v.Name.Name.String())
	}

	// 获取表主键列名,排除使用别名
	primaryKey := dbInstance.GetMeta().GetPrimaryKey(tableName)

	updateColumnsAndPrimaryKey := strings.Join(updateColumns, ",") + "," + primaryKey
	// 查询要更新字段数据的旧值，以及主键值
	selectSql := fmt.Sprintf("SELECT %s FROM %s %s LIMIT 200", updateColumnsAndPrimaryKey, tableStr, where)
	_, res, _ := dbInstance.SelectData(selectSql)

	bytes, _ := json.Marshal(res)
	dbSqlExec.OldValue = string(bytes)
	dbSqlExec.Table = tableName
	dbSqlExec.Type = entity.DbSqlExecTypeUpdate

	return doExec(execSqlReq.Sql, dbInstance)
}

func doDelete(delete *sqlparser.Delete, execSqlReq *DbSqlExecReq, dbSqlExec *entity.DbSqlExec) (*DbSqlExecRes, error) {
	dbInstance := execSqlReq.DbInstance

	tableStr := sqlparser.String(delete.TableExprs)
	// 可能使用别名，故空格切割
	table := strings.Split(tableStr, " ")[0]
	where := sqlparser.String(delete.Where)

	// 查询删除数据
	selectSql := fmt.Sprintf("SELECT * FROM %s %s LIMIT 200", tableStr, where)
	_, res, _ := dbInstance.SelectData(selectSql)

	bytes, _ := json.Marshal(res)
	dbSqlExec.OldValue = string(bytes)
	dbSqlExec.Table = table
	dbSqlExec.Type = entity.DbSqlExecTypeDelete

	return doExec(execSqlReq.Sql, dbInstance)
}

func doInsert(insert *sqlparser.Insert, execSqlReq *DbSqlExecReq, dbSqlExec *entity.DbSqlExec) (*DbSqlExecRes, error) {
	tableStr := sqlparser.String(insert.Table)
	// 可能使用别名，故空格切割
	table := strings.Split(tableStr, " ")[0]
	dbSqlExec.Table = table
	dbSqlExec.Type = entity.DbSqlExecTypeInsert

	return doExec(execSqlReq.Sql, execSqlReq.DbInstance)
}

func doExec(sql string, dbInstance *DbInstance) (*DbSqlExecRes, error) {
	rowsAffected, err := dbInstance.Exec(sql)
	if err != nil {
		return nil, err
	}
	res := make([]map[string]interface{}, 0)
	resData := make(map[string]interface{})
	resData["影响条数"] = rowsAffected
	res = append(res, resData)

	return &DbSqlExecRes{
		ColNames: []string{"影响条数"},
		Res:      res,
	}, nil
}
