package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

var cf config
var con *gorm.DB

func init() {
	flag.StringVar(&cf.FileName, "model-file", "", "Generate model file name")
	flag.StringVar(&cf.Directory, "model-directory", "", "Generated model directory")
	flag.StringVar(&cf.ModelName, "model-name", "", "Generate model struct name")
	flag.StringVar(&cf.DB, "connection", "", "DB connect dsn")
	flag.StringVar(&cf.TableName, "table", "", "Table name of generated model")
	flag.StringVar(&cf.DaoDirectory, "dao", "", "The directory of dao generate.")
	flag.StringVar(&cf.RepDirectory, "repo", "", "The directory of repository generate.")
	flag.StringVar(&cf.ConnDirectory, "conn", "", "The directory of conn generate.")
	flag.StringVar(&cf.ConfigFilePath, "config", "", "Special config file, format: .yml")
	flag.BoolVar(&cf.Force, "f", false, "force write")
}

func main() {
	flag.Parse()
	// load config
	if e := readConfigFromFile(&cf); e != nil {
		fmt.Println(e.Error())
		return
	}
	if e := cf.Validate(); e != nil {
		fmt.Println(e.Error())
		return
	}
	c, e := connect(cf.Getdsn())
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	con = c
	mp, e := getTableDescription()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	// if e := writeModelFile(mp); e != nil {
	// 	fmt.Println(e.Error())
	// 	return
	// }
	if e := writeDaoFile(mp); e != nil {
		fmt.Println(e.Error())
		return
	}
	if e := writeRepoFile(mp); e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println("Generate success!")
}

func getTableDescription() (*modelParse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}()
	tableName := cf.GetTableName()
	if con.Table(tableName).Statement.Table == "" {
		return nil, errors.New("table \"" + tableName + "\" not exist")
	}
	var result tableDcList
	con.Raw("DESCRIBE " + tableName).Scan(&result)
	modelDirectory := cf.GetDirectory()
	modelPackageName := "models"
	if len(modelDirectory) > 0 {
		sps := strings.Split(modelDirectory, "/")
		modelPackageName = sps[len(sps)-1]
	}
	daoPackageName := "dao"
	repoPackageName := "repo"
	if len(cf.DaoDirectory) > 0 {
		sps := strings.Split(cf.DaoDirectory, "/")
		daoPackageName = sps[len(sps)-1]
	}

	fmt.Println("ConnDirectory", cf.ConnDirectory)
	if len(cf.RepDirectory) > 0 {
		sps := strings.Split(cf.RepDirectory, "/")
		repoPackageName = sps[len(sps)-1]
	}
	Fields, Time := result.parseFields()
	parse := modelParse{
		Force:               cf.Force,
		ModelPackageName:    modelPackageName,
		ModelDirectory:      modelDirectory,
		ConnDirectory:       cf.ConnDirectory,
		FileName:            cf.GetFileName(),
		ModelName:           cf.GetModelName(),
		Fields:              Fields,
		TableName:           cf.GetTableName(),
		DaoDirectory:        cf.DaoDirectory,
		DaoPackageName:      daoPackageName,
		RepoPackageName:     repoPackageName,
		RepositoryDirectory: cf.RepDirectory,
		TimeImport:          Time,
	}

	return &parse, nil
}

func readConfigFromFile(cfg *config) error {
	configFile := cfg.ConfigFilePath
	// load default config if input config file is empty
	if len(configFile) == 0 {
		defaultFile, e := os.Stat(".yml")
		if e == nil && defaultFile.IsDir() == false {
			configFile = ".yml"
		}
	}
	if len(configFile) > 0 {
		bts, e := ioutil.ReadFile(configFile)
		if e != nil {
			return errors.New("config file open error:" + e.Error())
		}
		e = yaml.Unmarshal(bts, &cfg.fileConfig)
		// fmt.Println("fileConfig", g.Export(cfg.fileConfig))
		if e != nil {
			return errors.New("config file format error:" + e.Error())
		}
	}
	return nil
}
