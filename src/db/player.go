package db

type Player struct {
	PlayerId int64
	UserId int64
	PlayerName string
	Level int16
	Sex byte
	FightValue int32

	SceneId int16
	PosX int16
	PosY int16
	Direction byte
}