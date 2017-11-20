package main

import (
	"encoding/json"
	"io/ioutil"
)

type ManagerConfiguration struct {
	Admin            Admin
	Databases        []string
	UsersToIgnore    []User
	UsersPermissions []UserPermissions
}

type Admin struct {
	Name     string
	Password string
}

type User struct {
	Name string
	Host string
}

type UserPermissions struct {
	Name           string
	Host           string
	Password       string
	DatabaseGrants []DatabaseGrants
}

type DatabaseGrants struct {
	DatabaseName string
	Grants       []string
}

func getManagerConfiguration(path string) (*ManagerConfiguration, error) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var managerConfiguration ManagerConfiguration
	err = json.Unmarshal(jsonBytes, &managerConfiguration)

	return &managerConfiguration, err
}
