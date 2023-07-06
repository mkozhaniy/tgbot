package db

import (
	"database/sql"

	"github.com/tgbot/internal/domain/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Db     *gorm.DB
	Logger *zap.Logger
}

func Start(url string) (*gorm.DB, error) {
	sqlDb, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	if err != nil {
		return gormDb, err
	}
	return gormDb, nil
}

func InitTables(db *gorm.DB) error {
	if err := db.SetupJoinTable(&models.Bucket{}, "Products", &models.BucketProduct{}); err != nil {
		return err
	}
	if err := db.SetupJoinTable(&models.Order{}, "Products", &models.OrderProduct{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.Product{}, &models.Bucket{}); err != nil {
		return err
	}
	return nil
}

func DropTables(db *gorm.DB) error {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return err
	}
	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}
	return nil
}

func (pgdb Database) Save(ent interface{}) (interface{}, error) {
	switch ent.(type) {
	case models.Bucket:
		ent1 := ent.(models.Bucket)
		if err := models.PostgresHandleError(pgdb.Db.Save(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.Order:
		ent1 := ent.(models.Order)
		if err := models.PostgresHandleError(pgdb.Db.Save(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.User:
		ent1 := ent.(models.User)
		if err := models.PostgresHandleError(pgdb.Db.Save(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.Product:
		ent1 := ent.(models.Product)
		if err := models.PostgresHandleError(pgdb.Db.Save(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	}
	return nil, nil
}

func (pgdb Database) Delete(ent interface{}) (interface{}, error) {
	switch ent.(type) {
	case models.Bucket:
		ent1 := ent.(models.Bucket)
		if err := models.PostgresHandleError(pgdb.Db.
			Select(clause.Associations).Delete(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.Order:
		ent1 := ent.(models.Order)
		if err := models.PostgresHandleError(pgdb.Db.
			Select(clause.Associations).Delete(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.User:
		ent1 := ent.(models.User)
		if err := models.PostgresHandleError(pgdb.Db.
			Select(clause.Associations).Delete(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	case models.Product:
		ent1 := ent.(models.Product)
		if err := models.PostgresHandleError(pgdb.Db.
			Select(clause.Associations).Delete(&ent1).Error); err != nil {
			pgdb.Logger.Sugar().Info(err.Error())
			return nil, err
		}
		return ent1, nil
	}
	return nil, nil
}

func (pgdb Database) GetBucketByUserTgid(user_tgid int64) (*models.Bucket, error) {
	bucket := models.Bucket{UserTgid: user_tgid}
	err := pgdb.Db.Model(&models.Bucket{}).Preload("Products").First(&bucket).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return &bucket, nil
}
func (pgdb Database) GetOrderById(id uint) (*models.Order, error) {
	order := models.Order{ID: id}
	err := pgdb.Db.Model(&models.Order{}).Preload("Products").First(&order).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return &order, nil
}

func (pgdb Database) GetOrdersByUserTgid(user_tgid int64) ([]models.Order, error) {
	var orders []models.Order
	err := pgdb.Db.Model(&models.Order{}).Preload("Products").Find(&orders,
		"user_tgid = ?", user_tgid).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return orders, nil
}

func (pgdb Database) GetProductById(id uint) (*models.Product, error) {
	product := models.Product{ID: id}
	pgdb.Db.First(&product)
	if err := pgdb.Db.Error; err != nil {
		pgdb.Logger.Info(err.Error())
		return nil, err
	}
	return &product, nil

}
func (pgdb Database) GetPruductByName(name string) (*models.Product, error) {
	product := models.Product{Name: name}
	pgdb.Db.First(&product)
	if err := pgdb.Db.Error; err != nil {
		pgdb.Logger.Info(err.Error())
		return nil, err
	}
	return &product, nil
}

func (pgdb Database) GetUserByTgid(user_tgid int64) (*models.User, error) {
	user := models.User{Tgid: user_tgid}
	if err := pgdb.Db.First(&user).Error; err != nil {
		pgdb.Logger.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (pgdb Database) GetUserByName(name string) (*models.User, error) {
	user := models.User{Username: name}
	if err := pgdb.Db.First(&user).Error; err != nil {
		pgdb.Logger.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (pgdb Database) GetUserById(id uint) (*models.User, error) {
	user := models.User{ID: id}
	if err := pgdb.Db.First(&user).Error; err != nil {
		pgdb.Logger.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (pgdb Database) GetProductsByKind(kind string) ([]models.Product, error) {
	var result []models.Product
	err := pgdb.Db.Find(&result, "kind = ?", kind).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return result, nil
}

func (pgdb Database) GetAllProducts() ([]models.Product, error) {
	var result []models.Product
	err := pgdb.Db.Find(&result).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return result, nil
}

func (pgdb Database) GetProductByKind(kind string) (*models.Product, error) {
	var product models.Product
	err := pgdb.Db.First(&product, "kind = ?", kind).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return &product, nil
}

func (pgdb Database) GetAllKinds() ([]string, error) {
	var products []models.Product
	err := pgdb.Db.Distinct("kind").Order("kind").Find(&products).Error
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	kinds := make([]string, len(products))
	for i, val := range products {
		kinds[i] = val.Kind
	}
	return kinds, nil
}

func (pgdb Database) DeleteProductFromBucket(bucket *models.Bucket,
	product *models.Product) (*models.Bucket, error) {
	err := pgdb.Db.Model(bucket).Association("Products").Delete(product)
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return bucket, nil
}

func (pgdb Database) ClearBucket(bucket *models.Bucket) (*models.Bucket, error) {
	err := pgdb.Db.Model(bucket).Association("Products").Clear()
	if err != nil {
		pgdb.Logger.Sugar().Info(err.Error())
		return nil, err
	}
	return bucket, nil
}

func (pgdb Database) AddProductToBucket(bucket *models.Bucket,
	product *models.Product) (*models.Bucket, error) {
	err := pgdb.Db.Model(bucket).Association("Products").
		Append(product)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (pgdb Database) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := pgdb.Db.Model(&models.Order{}).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
