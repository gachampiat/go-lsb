package lsb


type Ilsb interface {
	InsertData(data []byte)(error)
	RetriveData()(msg []byte, err error)
	Detect()bool
}