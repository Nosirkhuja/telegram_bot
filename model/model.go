/*типы данных таблицы*/

package model

/*			  username	   ip   result    */
type Data map[string]map[string]string

type Hystory struct {
	IpId     string `gorm:"primaryKey;autoIncrement:false"`
	Ip       string
	Result   string
	Id       string
	Username string
}

type Admin struct {
	Username string
}
