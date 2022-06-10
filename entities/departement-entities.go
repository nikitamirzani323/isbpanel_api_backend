package entities

type Model_departement struct {
	Departement_id     string `json:"departement_id"`
	Departement_name   string `json:"departement_name"`
	Departement_create string `json:"departement_create"`
	Departement_update string `json:"departement_update"`
}
type Model_listdepartement struct {
	Departement_id   string `json:"departement_id"`
	Departement_name string `json:"departement_name"`
}

type Controller_departementsave struct {
	Page             string `json:"page" validate:"required"`
	Sdata            string `json:"sdata" validate:"required"`
	Departement_id   string `json:"departement_id" validate:"required"`
	Departement_name string `json:"departement_name" validate:"required"`
}
