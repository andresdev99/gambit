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
	CategID   int    `json:"categID" db:"Categ_Id"`
	CategName string `json:"categName" db:"Categ_Name"`
	CategPath string `json:"categPath" db:"Categ_Path"`
}

type Product struct {
	ProdID          int     `json:"prodID" db:"Prod_Id"`
	ProdTitle       string  `json:"prodTitle" db:"Prod_Title"`
	ProdDescription string  `json:"prodDescription" db:"Prod_Description"`
	ProdCreatedAt   string  `json:"prodCreatedAt" db:"Prod_CreatedAt"`
	ProdUpdated     string  `json:"prodUpdated" db:"Prod_Updated"`
	ProdPrice       float64 `json:"prodPrice,omitempty" db:"Prod_Price"`
	ProdPath        string  `json:"prodPath" db:"Prod_Path"`
	ProdCategoryID  int     `json:"prodCategoryId" db:"Prod_CategoryId"`
	ProdStock       int     `json:"prodStock" db:"Prod_Stock"`

	// These fields don't exist in the database, so they shouldn't have `db` tags
	ProdSearch    int `json:"prodSearch,omitempty"`
	ProdCategPath int `json:"ProdCategPath,omitempty"`
}
