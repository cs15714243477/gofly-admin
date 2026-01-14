package dbtablemgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gctx"
)

// 数据库表管理（建议仅开发环境使用）
type Index struct {
	CustomRoutes gf.Map
}

func init() {
	fpath := Index{
		CustomRoutes: gf.Map{
			// 默认自动路由不会把 Drop* 识别为 DELETE，这里显式声明
			"DropTable":  "DELETE:/business/dbtablemgr/dropTable",
			"DropColumn": "DELETE:/business/dbtablemgr/dropColumn",
		},
	}
	gf.Register(&fpath, fpath)
}

var (
	ctx = gctx.New()
)

func mustDevEnv() error {
	appConf, _ := gcfg.Instance().Get(ctx, "app")
	runEnv := gconv.String(gconv.Map(appConf)["runEnv"])
	if runEnv == "release" {
		return errors.New("生产环境禁止操作")
	}
	return nil
}

var identRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,63}$`)

func validateIdent(s string) bool {
	return identRe.MatchString(s)
}

func quoteIdent(name string) string {
	// name 必须先经过 validateIdent
	return "`" + name + "`"
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

type CreateTableField struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Length        string      `json:"length"` // 允许 "64" 或 "10,2"
	Unsigned      bool        `json:"unsigned"`
	Nullable      bool        `json:"nullable"`
	Default       interface{} `json:"default"` // nil 表示不设置默认值；"" 表示默认空字符串
	Comment       string      `json:"comment"`
	PrimaryKey    bool        `json:"primaryKey"`
	AutoIncrement bool        `json:"autoIncrement"`
}

type CreateTableReq struct {
	TableName    string             `json:"tableName"`
	TableComment string             `json:"tableComment"`
	Engine       string             `json:"engine"`
	Charset      string             `json:"charset"`
	Collate      string             `json:"collate"`
	Fields       []CreateTableField `json:"fields"`
}

type AlterColumnReq struct {
	TableName string           `json:"tableName"`
	OldName   string           `json:"oldName"`  // 修改/重命名时使用
	NewName   string           `json:"newName"`  // 可选：重命名
	Field     CreateTableField `json:"field"`    // 字段定义（name/type/length/nullable/default/comment/unsigned/primaryKey/autoIncrement）
	Position  string           `json:"position"` // 可选：FIRST / AFTER
	After     string           `json:"after"`    // 可选：当 Position=AFTER 时指定字段名
}

var (
	allowedTypes = map[string]struct{}{
		"bool":    {},
		"tinyint": {}, "smallint": {}, "int": {}, "bigint": {},
		"float": {}, "double": {}, "decimal": {},
		"char": {}, "varchar": {}, "text": {}, "tinytext": {}, "mediumtext": {}, "longtext": {},
		"date": {}, "datetime": {}, "timestamp": {},
		"blob": {}, "json": {},
	}
	typesNeedLength = map[string]struct{}{
		"char": {}, "varchar": {}, "decimal": {},
	}
)

func isNumericType(t string) bool {
	switch t {
	case "bool":
		return true
	case "tinyint", "smallint", "int", "bigint", "float", "double", "decimal":
		return true
	default:
		return false
	}
}

func buildColumnType(t, length string, unsigned bool) (string, error) {
	tt := strings.ToLower(strings.TrimSpace(t))
	// bool 统一映射为 tinyint(1)
	if tt == "bool" {
		tt = "tinyint"
		if strings.TrimSpace(length) == "" {
			length = "1"
		}
	}
	if _, ok := allowedTypes[tt]; !ok {
		return "", errors.New("不支持的字段类型: " + t)
	}
	length = strings.TrimSpace(length)
	if _, need := typesNeedLength[tt]; need && length == "" {
		return "", errors.New("字段类型 " + tt + " 需要 length")
	}
	colType := strings.ToUpper(tt)
	if length != "" {
		// 对 length 做最小约束，避免注入：只允许数字、逗号
		for _, ch := range length {
			if (ch < '0' || ch > '9') && ch != ',' {
				return "", errors.New("length 不合法，仅允许数字或逗号")
			}
		}
		colType += fmt.Sprintf("(%s)", length)
	}
	if unsigned && (tt == "tinyint" || tt == "smallint" || tt == "int" || tt == "bigint" || tt == "decimal") {
		colType += " UNSIGNED"
	}
	return colType, nil
}

func buildDefaultClause(field CreateTableField) (string, error) {
	if field.Default == nil {
		return "", nil
	}
	tt := strings.ToLower(strings.TrimSpace(field.Type))
	// 允许 CURRENT_TIMESTAMP（不加引号）
	defStr := strings.TrimSpace(gconv.String(field.Default))
	if defStr == "" && field.Default != "" {
		// field.Default 是空字符串时 gconv.String 也是空字符串；这里仅防御 nil/未知类型
	}
	up := strings.ToUpper(defStr)
	if (tt == "timestamp" || tt == "datetime") && (up == "CURRENT_TIMESTAMP") {
		return " DEFAULT " + up, nil
	}
	if isNumericType(tt) {
		// 数值类型：尽量不加引号（基础校验：允许数字/小数点/负号）
		for i, ch := range defStr {
			if (ch < '0' || ch > '9') && ch != '.' && ch != '-' {
				if i == 0 && ch == '+' {
					continue
				}
				return "", errors.New("数值类型默认值不合法")
			}
		}
		return " DEFAULT " + defStr, nil
	}
	// 其他类型：字符串默认值，做单引号转义
	return " DEFAULT '" + escapeSingleQuotes(defStr) + "'", nil
}

func buildCreateTableSQL(req CreateTableReq) (string, error) {
	if req.TableName == "" || !validateIdent(req.TableName) {
		return "", errors.New("tableName 不合法")
	}
	if len(req.Fields) == 0 {
		return "", errors.New("fields 不能为空")
	}

	engine := strings.TrimSpace(req.Engine)
	if engine == "" {
		engine = "InnoDB"
	}
	charset := strings.TrimSpace(req.Charset)
	if charset == "" {
		charset = "utf8mb4"
	}
	collate := strings.TrimSpace(req.Collate)

	lines := make([]string, 0, len(req.Fields)+2)
	pks := make([]string, 0)
	seen := make(map[string]struct{}, len(req.Fields))

	for _, f := range req.Fields {
		f.Name = strings.TrimSpace(f.Name)
		if f.Name == "" || !validateIdent(f.Name) {
			return "", errors.New("字段名不合法: " + f.Name)
		}
		if _, ok := seen[f.Name]; ok {
			return "", errors.New("字段名重复: " + f.Name)
		}
		seen[f.Name] = struct{}{}

		colType, err := buildColumnType(f.Type, f.Length, f.Unsigned)
		if err != nil {
			return "", err
		}

		nullClause := " NOT NULL"
		if f.Nullable && !f.PrimaryKey && !f.AutoIncrement {
			nullClause = " NULL"
		}
		defClause, err := buildDefaultClause(f)
		if err != nil {
			return "", err
		}
		extra := ""
		if f.AutoIncrement {
			extra += " AUTO_INCREMENT"
		}
		commentClause := ""
		if strings.TrimSpace(f.Comment) != "" {
			commentClause = " COMMENT '" + escapeSingleQuotes(strings.TrimSpace(f.Comment)) + "'"
		}

		line := fmt.Sprintf("  %s %s%s%s%s", quoteIdent(f.Name), colType, nullClause, defClause, extra)
		line += commentClause
		lines = append(lines, line)

		if f.PrimaryKey {
			pks = append(pks, quoteIdent(f.Name))
		}
	}

	if len(pks) > 0 {
		sort.Strings(pks)
		lines = append(lines, "  PRIMARY KEY ("+strings.Join(pks, ", ")+")")
	}

	sql := "CREATE TABLE " + quoteIdent(req.TableName) + " (\n" + strings.Join(lines, ",\n") + "\n)"
	sql += " ENGINE=" + engine
	sql += " DEFAULT CHARSET=" + charset
	if collate != "" {
		sql += " COLLATE=" + collate
	}
	if strings.TrimSpace(req.TableComment) != "" {
		sql += " COMMENT='" + escapeSingleQuotes(strings.TrimSpace(req.TableComment)) + "'"
	}
	sql += ";"
	return sql, nil
}

func buildColumnDefSQL(name string, f CreateTableField) (string, error) {
	if name == "" || !validateIdent(name) {
		return "", errors.New("字段名不合法")
	}
	colType, err := buildColumnType(f.Type, f.Length, f.Unsigned)
	if err != nil {
		return "", err
	}
	nullClause := " NOT NULL"
	if f.Nullable && !f.PrimaryKey && !f.AutoIncrement {
		nullClause = " NULL"
	}
	defClause, err := buildDefaultClause(f)
	if err != nil {
		return "", err
	}
	extra := ""
	if f.AutoIncrement {
		extra += " AUTO_INCREMENT"
	}
	commentClause := ""
	if strings.TrimSpace(f.Comment) != "" {
		commentClause = " COMMENT '" + escapeSingleQuotes(strings.TrimSpace(f.Comment)) + "'"
	}
	return fmt.Sprintf("%s %s%s%s%s%s", quoteIdent(name), colType, nullClause, defClause, extra, commentClause), nil
}

func buildAlterPositionSQL(pos, after string) (string, error) {
	pos = strings.ToUpper(strings.TrimSpace(pos))
	if pos == "" {
		return "", nil
	}
	if pos == "FIRST" {
		return " FIRST", nil
	}
	if pos == "AFTER" {
		after = strings.TrimSpace(after)
		if after == "" || !validateIdent(after) {
			return "", errors.New("after 字段名不合法")
		}
		return " AFTER " + quoteIdent(after), nil
	}
	return "", errors.New("position 仅支持 FIRST 或 AFTER")
}

// 获取数据表（GET /business/dbtablemgr/getTables）
func (api *Index) GetTables(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	dbConf, _ := gcfg.Instance().Get(ctx, "database.default")
	dbname := gconv.String(gconv.Map(dbConf)["dbname"])
	tablelist, err := gf.DB().Query(c, "select TABLE_NAME,TABLE_COMMENT from information_schema.tables where table_schema = '"+dbname+"'")
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	var out []interface{}
	for _, v := range tablelist {
		out = append(out, map[string]interface{}{"name": v["TABLE_NAME"], "title": v["TABLE_COMMENT"]})
	}
	gf.Success().SetMsg("获取数据表").SetData(out).Regin(c)
}

// 获取字段（GET /business/dbtablemgr/getColumns?tablename=xxx）
func (api *Index) GetColumns(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	tablename := c.DefaultQuery("tablename", "")
	if tablename == "" || !validateIdent(tablename) {
		gf.Failed().SetMsg("tablename不合法").Regin(c)
		return
	}
	dbConf, _ := gcfg.Instance().Get(ctx, "database.default")
	dbname := gconv.String(gconv.Map(dbConf)["dbname"])
	rows, err := gf.DB().Query(c,
		"select COLUMN_NAME,COLUMN_COMMENT,COLUMN_TYPE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,IS_NULLABLE,COLUMN_DEFAULT,NUMERIC_PRECISION,EXTRA "+
			"from information_schema.columns where TABLE_SCHEMA='"+dbname+"' AND TABLE_NAME='"+tablename+"' order by ORDINAL_POSITION",
	)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取字段").SetData(rows).Regin(c)
}

// 获取建表SQL（GET /business/dbtablemgr/getCreateSQL?tablename=xxx）
func (api *Index) GetCreateSQL(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	tablename := c.DefaultQuery("tablename", "")
	if tablename == "" || !validateIdent(tablename) {
		gf.Failed().SetMsg("tablename不合法").Regin(c)
		return
	}
	rows, err := gf.DB().Query(c, "SHOW CREATE TABLE `"+tablename+"`")
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取建表SQL").SetData(rows).Regin(c)
}

// 预览建表SQL（POST /business/dbtablemgr/previewCreateSQL）
// body: { tableName, tableComment, engine, charset, collate, fields: [...] }
func (api *Index) PreviewCreateSQL(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	raw, _ := json.Marshal(param)
	var req CreateTableReq
	if uerr := json.Unmarshal(raw, &req); uerr != nil {
		gf.Failed().SetMsg("参数解析失败").SetData(uerr.Error()).Regin(c)
		return
	}
	sqlStr, berr := buildCreateTableSQL(req)
	if berr != nil {
		gf.Failed().SetMsg("生成建表SQL失败").SetData(berr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("生成建表SQL成功").SetData(gf.Map{"sql": sqlStr}).Regin(c)
}

// 可视化创建表（POST /business/dbtablemgr/createTable）
// body: 同 PreviewCreateSQL
func (api *Index) CreateTable(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	raw, _ := json.Marshal(param)
	var req CreateTableReq
	if uerr := json.Unmarshal(raw, &req); uerr != nil {
		gf.Failed().SetMsg("参数解析失败").SetData(uerr.Error()).Regin(c)
		return
	}
	sqlStr, berr := buildCreateTableSQL(req)
	if berr != nil {
		gf.Failed().SetMsg("生成建表SQL失败").SetData(berr.Error()).Regin(c)
		return
	}
	_, exErr := gf.DB().Exec(c, sqlStr)
	if exErr != nil {
		gf.Failed().SetMsg("创建失败").SetData(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("创建成功").SetData(gf.Map{"sql": sqlStr}).Regin(c)
}

// 新增字段（POST /business/dbtablemgr/addColumn）
func (api *Index) AddColumn(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	raw, _ := json.Marshal(param)
	var req AlterColumnReq
	if uerr := json.Unmarshal(raw, &req); uerr != nil {
		gf.Failed().SetMsg("参数解析失败").SetData(uerr.Error()).Regin(c)
		return
	}
	req.TableName = strings.TrimSpace(req.TableName)
	if req.TableName == "" || !validateIdent(req.TableName) {
		gf.Failed().SetMsg("tableName 不合法").Regin(c)
		return
	}
	req.Field.Name = strings.TrimSpace(req.Field.Name)
	if req.Field.Name == "" || !validateIdent(req.Field.Name) {
		gf.Failed().SetMsg("字段名不合法").Regin(c)
		return
	}
	colDef, derr := buildColumnDefSQL(req.Field.Name, req.Field)
	if derr != nil {
		gf.Failed().SetMsg("生成字段SQL失败").SetData(derr.Error()).Regin(c)
		return
	}
	posSQL, perr := buildAlterPositionSQL(req.Position, req.After)
	if perr != nil {
		gf.Failed().SetMsg("字段位置参数错误").SetData(perr.Error()).Regin(c)
		return
	}
	sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s%s;", quoteIdent(req.TableName), colDef, posSQL)
	_, exErr := gf.DB().Exec(c, sqlStr)
	if exErr != nil {
		gf.Failed().SetMsg("新增字段失败").SetData(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("新增字段成功").SetData(gf.Map{"sql": sqlStr}).Regin(c)
}

// 修改/重命名字段（POST /business/dbtablemgr/modifyColumn）
func (api *Index) ModifyColumn(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	raw, _ := json.Marshal(param)
	var req AlterColumnReq
	if uerr := json.Unmarshal(raw, &req); uerr != nil {
		gf.Failed().SetMsg("参数解析失败").SetData(uerr.Error()).Regin(c)
		return
	}
	req.TableName = strings.TrimSpace(req.TableName)
	if req.TableName == "" || !validateIdent(req.TableName) {
		gf.Failed().SetMsg("tableName 不合法").Regin(c)
		return
	}
	req.OldName = strings.TrimSpace(req.OldName)
	if req.OldName == "" || !validateIdent(req.OldName) {
		gf.Failed().SetMsg("oldName 不合法").Regin(c)
		return
	}
	newName := strings.TrimSpace(req.NewName)
	if newName == "" {
		newName = req.OldName
	}
	if !validateIdent(newName) {
		gf.Failed().SetMsg("newName 不合法").Regin(c)
		return
	}
	// 字段定义里的 name 以 newName 为准
	req.Field.Name = newName
	colDef, derr := buildColumnDefSQL(newName, req.Field)
	if derr != nil {
		gf.Failed().SetMsg("生成字段SQL失败").SetData(derr.Error()).Regin(c)
		return
	}
	posSQL, perr := buildAlterPositionSQL(req.Position, req.After)
	if perr != nil {
		gf.Failed().SetMsg("字段位置参数错误").SetData(perr.Error()).Regin(c)
		return
	}
	sqlStr := ""
	if newName != req.OldName {
		// CHANGE 支持改名 + 修改定义
		sqlStr = fmt.Sprintf("ALTER TABLE %s CHANGE COLUMN %s %s%s;", quoteIdent(req.TableName), quoteIdent(req.OldName), colDef, posSQL)
	} else {
		sqlStr = fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s%s;", quoteIdent(req.TableName), colDef, posSQL)
	}
	_, exErr := gf.DB().Exec(c, sqlStr)
	if exErr != nil {
		gf.Failed().SetMsg("修改字段失败").SetData(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("修改字段成功").SetData(gf.Map{"sql": sqlStr}).Regin(c)
}

// 删除字段（DELETE /business/dbtablemgr/dropColumn）
func (api *Index) DropColumn(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	tableName := strings.TrimSpace(gconv.String(param["tableName"]))
	colName := strings.TrimSpace(gconv.String(param["name"]))
	if tableName == "" || !validateIdent(tableName) {
		gf.Failed().SetMsg("tableName 不合法").Regin(c)
		return
	}
	if colName == "" || !validateIdent(colName) {
		gf.Failed().SetMsg("字段名不合法").Regin(c)
		return
	}
	sqlStr := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", quoteIdent(tableName), quoteIdent(colName))
	_, exErr := gf.DB().Exec(c, sqlStr)
	if exErr != nil {
		gf.Failed().SetMsg("删除字段失败").SetData(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("删除字段成功").SetData(gf.Map{"sql": sqlStr}).Regin(c)
}

// 创建表（POST /business/dbtablemgr/createTableBySQL）
// body: { "sql": "CREATE TABLE ..." }
func (api *Index) CreateTableBySQL(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	sqlStr := strings.TrimSpace(gconv.String(param["sql"]))
	if sqlStr == "" {
		gf.Failed().SetMsg("sql不能为空").Regin(c)
		return
	}
	// 禁止多语句
	if strings.Contains(sqlStr, ";") {
		gf.Failed().SetMsg("不允许包含分号（禁止多语句）").Regin(c)
		return
	}
	up := strings.ToUpper(strings.TrimSpace(sqlStr))
	if !strings.HasPrefix(up, "CREATE TABLE") {
		gf.Failed().SetMsg("仅允许CREATE TABLE语句").Regin(c)
		return
	}
	if strings.Contains(up, "DROP ") || strings.Contains(up, "DELETE ") || strings.Contains(up, "UPDATE ") || strings.Contains(up, "INSERT ") {
		gf.Failed().SetMsg("SQL包含危险关键字").Regin(c)
		return
	}
	_, exErr := gf.DB().Exec(c, sqlStr)
	if exErr != nil {
		gf.Failed().SetMsg(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("创建成功").Regin(c)
}

// 删除表（DELETE /business/dbtablemgr/dropTable）
// body: { "tablename": "xxx" }
func (api *Index) DropTable(c *gf.GinCtx) {
	if err := mustDevEnv(); err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	param, err := gf.RequestParam(c)
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	tablename := strings.TrimSpace(gconv.String(param["tablename"]))
	if tablename == "" || !validateIdent(tablename) {
		gf.Failed().SetMsg("tablename不合法").Regin(c)
		return
	}
	_, exErr := gf.DB().Exec(c, "DROP TABLE IF EXISTS `"+tablename+"`")
	if exErr != nil {
		gf.Failed().SetMsg(exErr.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("删除成功").Regin(c)
}
