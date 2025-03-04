package main

import (
	"fmt"
	"os"
	"path"
)

type parseField struct {
	Attr       string
	Type       fieldType
	Tag        string
	ColumnName string
	IsPrimary  bool
}

type modelParse struct {
	ModelPackageName    string
	DaoPackageName      string
	RepoPackageName     string
	FileName            string
	ModelName           string
	TimeImport          bool
	Fields              []parseField
	TableName           string
	ConnDirectory       string
	ModelDirectory      string
	RepositoryDirectory string
	DaoDirectory        string
	Force               bool
}

func (m modelParse) mysqlDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	fmt.Println("m.ConnDirectory", m.ConnDirectory)
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + m.ConnDirectory
	}
	return m.ConnDirectory
}

func (m modelParse) daoDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + m.DaoDirectory
	}
	return m.DaoDirectory
}

func (m modelParse) modelDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + m.ModelDirectory
	}
	return m.ModelDirectory
}

func (m modelParse) primaryKey() string {
	for _, value := range m.Fields {
		if value.IsPrimary {
			return value.ColumnName
		}
	}
	return "id"
}

func (m modelParse) primaryKeyType() fieldType {
	for _, value := range m.Fields {
		if value.IsPrimary {
			return value.Type
		}
	}
	return TypeUnknown
}

func (m modelParse) RepositoryInterfaceName() string {
	return m.ModelName + "Repository"
}

func (m modelParse) DaoStructName() string {
	return m.ModelName + "Dao"
}
