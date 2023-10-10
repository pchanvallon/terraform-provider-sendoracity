package provider

type City struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Touristic *bool  `json:"touristic,omitempty"`
}

type House struct {
	Id          int    `json:"id,omitempty"`
	CityId      int    `json:"cityid,omitempty"`
	Address     string `json:"address,omitempty"`
	Inhabitants int    `json:"inhabitants,omitempty"`
}

type Store struct {
	Id      int    `json:"id,omitempty"`
	CityId  int    `json:"cityid,omitempty"`
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
}
