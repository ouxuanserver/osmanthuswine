package helper

import (
	"github.com/go-errors/errors"
	"github.com/wailovet/db2struct"
	"github.com/wailovet/osmanthuswine/src/core"
	"strconv"
)

func GetStructByDb(tableName string, packageName string, structName string) (string, error) {
	mariadbUser := core.GetInstanceConfig().Db.User
	mariadbPassword := core.GetInstanceConfig().Db.Password
	mariadbHost := core.GetInstanceConfig().Db.Host
	mariadbPort, _ := strconv.Atoi(core.GetInstanceConfig().Db.Port)
	mariadbDatabase := core.GetInstanceConfig().Db.Name
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase, tableName)

	if err != nil {
		return "", errors.New("Error in selecting column data information from mysql information schema")
	}

	// Generate struct string based on columnDataTypes
	struc, err := db2struct.Generate(*columnDataTypes, tableName, structName, packageName, true, true, false)

	if err != nil {
		return "", errors.New("Error in creating struct from json: " + err.Error())
	}
	return string(struc), nil

}
