package repositories

type Repository struct {
	CatRepository CatRepository
}

func SetupRepository() *Repository {
	return &Repository{
		CatRepository: NewCatRepository(),
	}
}
