package account

type Handler struct {
	accountSvc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{accountSvc: svc}
}
