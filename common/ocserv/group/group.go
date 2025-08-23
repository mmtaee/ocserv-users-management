package group

type OcservGroup struct{}

type OcservGroupInterface interface{}

func NewOcservGroup() *OcservGroup {
	return &OcservGroup{}
}
