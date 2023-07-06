package interfaces

type EntityStorage interface {
	Save(ent interface{}) (interface{}, error)
	Delete(ent interface{}) (interface{}, error)
}
