package models

type SecretRDSJson struct {
	UserName            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
}

type SignUp struct {
	UserEmail string `json:"UserEmail"`
	UserUUID  string `json:"UserUUID"`
}

type Category struct {
	CategID   int    `json:"categID"`
	CategName string `json:"categName"`
	CategPath string `json:"categPath"`
}

type Product struct {
	ProdID          int     `json:"prodID"`
	ProdTitle       string  `json:"prodTitle"`
	ProdDescription string  `json:"prodDescription"`
	ProdCreatedAt   string  `json:"prodCreatedAt"`
	ProdUpdated     string  `json:"prodUpdated"`
	ProdPrice       float64 `json:"prodPrice,omitempty"`
	ProdPath        string  `json:"prodPath"`
	ProdCategoryID  int     `json:"prodCategoryId"`
	ProdStock       int     `json:"prodStock"`
	ProdSearch      int     `json:"prodSearch,omitempty"`
	ProdCategPath   int     `json:"ProdCategPath,omitempty"`
}
