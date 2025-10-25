package support

type SupportService struct {
	supportRepo *SupportRepo
}

func NewSupportService(supportRepo *SupportRepo) *SupportService {
	return &SupportService{
		supportRepo: supportRepo,
	}
}
