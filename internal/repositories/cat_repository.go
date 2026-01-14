package repositories

type CatRepository interface {
	GetCatDB() (string, error)
}

type catRepository struct {
}

func NewCatRepository() CatRepository {
	return &catRepository{}
}

func (cr *catRepository) GetCatDB() (string, error) {
	return "meow", nil
}
