package model

type API struct {
	Name      string      `json:"name"`
	Portal    string      `json:"portal"`
	Resources []*Resource `json:"resources"`
	Key       []byte
}

func (a *API) FindResourceByPath(p string) *Resource {
	for _, v := range a.Resources {
		if v.Path == p {
			return v
		}
	}
	return nil
}
