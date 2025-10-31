module github.com/mhshajib/chronz/examples/gorm_postgres

go 1.22

require (
	github.com/mhshajib/chronz v0.0.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
)

replace github.com/mhshajib/chronz => ../../
