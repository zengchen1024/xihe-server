package bigmodels

type Config struct {
	User         string `json:"user"            required:"true"`
	Password     string `json:"password"        required:"true"`
	Project      string `json:"project"         required:"true"`
	AuthEndpoint string `json:"auth_endpoing"   required:"true"`

	EndpointOfDescribingPicture string `json:"endpoint_of_describing_picture" required:"true"`
}
