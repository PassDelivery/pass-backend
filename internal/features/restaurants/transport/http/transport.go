package restaurants_transport_http


type RestaurantsHTTPHandler struct {
	restaurantsService RestaurantsService
}


type RestaurantsService interface {

}

func NewRestaurantsHTTPHandler(rs RestaurantsService) *RestaurantsHTTPHandler{
	return &RestaurantsHTTPHandler{
		restaurantsService: rs,
	}
}

