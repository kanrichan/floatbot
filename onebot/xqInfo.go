package onebot

import (
	"reflect"
	"yaya/core"

	"github.com/tidwall/gjson"
)

type XGroupInfo struct {
	GroupID        int64  `db:"group_id" json:"group_id"`
	GroupName      string `db:"group_name" json:"group_name"`
	MemberCount    int64  `db:"member_count" json:"member_count"`
	MaxMemberCount int64  `db:"max_member_count" json:"max_member_count"`
}

type XGroupMember struct {
	GroupID         int64  `db:"group_id" json:"group_id"`
	UserID          int64  `db:"user_id" json:"user_id"`
	Nickname        string `db:"nickname" json:"nickname"`
	Card            string `db:"card" json:"card"`
	Sex             string `db:"sex" json:"sex"`
	Age             int64  `db:"age" json:"age"`
	Area            string `db:"area" json:"area"`
	JoinTime        int64  `db:"join_time" json:"join_time"`
	LastSentTime    int64  `db:"last_sent_time" json:"last_sent_time"`
	Level           string `db:"level" json:"level"`
	Role            string `db:"role" json:"role"`
	Unfriendly      bool   `db:"unfriendly" json:"unfriendly"`
	Title           string `db:"title" json:"title"`
	TitleExpireTime int64  `db:"title_expire_time" json:"title_expire_time"`
	CardChangeable  bool   `db:"card_changeable" json:"card_changeable"`
}

func (bot *BotYaml) saveGroupInfo() {
	groupList := core.GetGroupList(bot.Bot)
	if groupList == "" {
		return
	}

	g := gjson.Parse(groupList)
	for _, o := range append(g.Get("create").Array(), append(g.Get("manage").Array(), g.Get("join").Array()...)...) {
		memberList := core.GetGroupMemberList_B(bot.Bot, o.Get("gc").Int())
		m := gjson.Parse(memberList)
		info := XGroupInfo{
			GroupID:        o.Get("gc").Int(),
			GroupName:      unicode2chinese(o.Get("gn").Str),
			MemberCount:    m.Get("mem_num").Int(),
			MaxMemberCount: m.Get("max_num").Int(),
		}
		bot.dbInsert(info)

		membersMap := m.Get("members").Map()
		list := reflect.ValueOf(membersMap).MapKeys()
		for _, member := range list {

			qq := member.Interface().(string)
			nickname := m.Get("members." + qq + ".nk").Str
			card := m.Get("members." + qq + ".cd").Str
			role := "member"
			for _, admin := range m.Get("adm").Array() {
				if qq == admin.Str {
					role = "admin"
				}
			}
			if qq == m.Get("owner").Str {
				role = "owner"
			}
			member := XGroupMember{
				GroupID:         o.Get("gc").Int(),
				UserID:          core.Str2Int(qq),
				Nickname:        nickname,
				Card:            card,
				Sex:             "unknown",
				Age:             0,
				Area:            "",
				JoinTime:        m.Get("members." + qq + ".jt").Int(),
				LastSentTime:    m.Get("members." + qq + ".lst").Int(),
				Level:           m.Get("members." + qq + ".ll").Str,
				Role:            role,
				Unfriendly:      false,
				Title:           "",
				TitleExpireTime: 0,
				CardChangeable:  false,
			}
			bot.dbInsert(member)
		}
	}
}
