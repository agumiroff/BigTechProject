package model

type PartRow struct {
	PartUUID  string `db:"part_uuid"`
	OrderUUID string `db:"order_uuid"`
}
