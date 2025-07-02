package models

type Order struct {
	UserUUID  string   `json:"user_uuid"`
	PartUUIDs []string `json:"part_uuids"`
}
