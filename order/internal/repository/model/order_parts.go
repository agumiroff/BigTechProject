package model

type OrderParts struct {
	OrderUUID string `db:"order_uuid"`
	PartUUID  string `db:"part_uuid"`
}
