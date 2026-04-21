package user

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
	FindAll() ([]User, error)
	FindAllWithFilter(page, limit int, search, role string) ([]User, int64, error)
	Delete(id uint) error
}

type GormRepository struct {
	db *gorm.DB
}

var _ Repository = (*GormRepository)(nil)

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *GormRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *GormRepository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *GormRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *GormRepository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *GormRepository) FindAllWithFilter(page, limit int, search, role string) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.Model(&User{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if role != "" {
		query = query.Where("role = ? AND role != ?", role, "admin")
	} else {
		query = query.Where("role != ?", "admin")
	}

	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *GormRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
