package database

import "errors"

// ErrNotFound is returned when a requested record does not exist in the
// database. Handlers can check for it with errors.Is to return a proper HTTP
// 404 response.
var ErrNotFound = errors.New("registro não encontrado")

var ErrGenericDatabase = errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")

var ErrProcessQuery = errors.New("Falha ao processar os resultados da consulta, em caso de persistencia contatar o suporte.")

var ErrBindQuery = errors.New("Falha indexar os resultados, em caso de persistencia contatar o suporte.")

func ErrRegisterObject(object string) error {
	return errors.New("Falha ao registrar um(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}

func ErrUpdateObject(object string) error {
	return errors.New("Falha ao atualizar o(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}

func ErrDeleteObject(object string) error {
	return errors.New("Falha ao deletar o(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}
