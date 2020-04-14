package lsb


type Ilsb interface {
	InsertData(data []byte)(error)
	RetriveData(lenght int)(msg []byte, err error)
}