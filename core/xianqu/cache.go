package xianqu

import "C"

import (
	"reflect"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

type CacheData struct {
	M     sync.Mutex
	Max   int
	Key   []interface{}
	Value []interface{}
}

func (c *CacheData) Insert(key interface{}, value interface{}) {
	c.M.Lock()
	defer c.M.Unlock()
	switch {
	case len(c.Key) >= c.Max:
		start := (c.Max + 10) / 10
		c.Key = c.Key[start:]
		c.Value = c.Value[start:]
		fallthrough
	default:
		c.Key = append(c.Key, key)
		c.Value = append(c.Value, value)
	}
}

func (c *CacheData) Search(key interface{}) (value interface{}) {
	length := len(c.Key)
	for i := length - 1; i >= 0; i-- {
		if key == c.Key[i] {
			return c.Value[i]
		}
	}
	return nil
}

// Hcraes Value反向搜索Key
func (c *CacheData) Hcraes(value interface{}) (key interface{}) {
	length := len(c.Value)
	for i := length - 1; i >= 0; i-- {
		if value == c.Value[i] {
			return c.Key[i]
		}
	}
	return nil
}

type CacheGroupsData struct {
	M     sync.Mutex
	Group []*GroupData
}

type GroupData struct {
	GroupInfo    *GroupInfo
	GroupMembers []*GroupMember
}

type GroupInfo struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int64  `json:"member_count"`
	MaxMemberCount int64  `json:"max_member_count"`
}

type GroupMember struct {
	GroupID         int64  `json:"group_id"`
	UserID          int64  `json:"user_id"`
	Nickname        string `json:"nickname"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int64  `json:"age"`
	Area            string `json:"area"`
	JoinTime        int64  `json:"join_time"`
	LastSentTime    int64  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int64  `json:"title_expire_time"`
	CardChangeable  bool   `json:"card_changeable"`
}

// GetCacheGroup 获取群信息，不存在则请求
func (c *CacheGroupsData) GetCacheGroup(bot, groupID int64, cache bool) (group *GroupData) {
	if cache { // 使用缓存
		for i := range c.Group {
			if groupID == c.Group[i].GroupInfo.GroupID {
				return c.Group[i] // 存在即立刻返回数据
			}
		}
	}
	group = &GroupData{}
	// 向服务器请求
	name := XQApiGroupName(bot, groupID)
	listB := XQApiGroupMemberListB(bot, groupID)
	if listB != "" {
		listGJson := gjson.Parse(listB)
		membersGJson := listGJson.Get("members")
		qqReflectValues := reflect.ValueOf(membersGJson.Map()).MapKeys()
		owner := listGJson.Get("owner").Int()
		admins := listGJson.Get("adm").Array()
		group.GroupInfo = &GroupInfo{
			GroupID:        groupID,
			GroupName:      name,
			MemberCount:    listGJson.Get("mem_num").Int(),
			MaxMemberCount: listGJson.Get("max_num").Int(),
		}
		for _, value := range qqReflectValues {
			qq := value.Interface().(string)
			role := "member"
			for i := range admins {
				if str2Int(qq) == admins[i].Int() {
					role = "admin"
				}
			}
			if qq == int2Str(owner) {
				role = "owner"
			}
			nickname := strings.ReplaceAll(membersGJson.Get(qq).Get("nk").Str, "&nbsp;", " ")
			card := strings.ReplaceAll(membersGJson.Get(qq).Get("cd").Str, "&nbsp;", " ")
			if card == "" {
				card = nickname
			}
			member := &GroupMember{
				GroupID:         groupID,
				UserID:          str2Int(qq),
				Nickname:        nickname,
				Card:            card,
				Sex:             "",
				Age:             0,
				Area:            "unknown",
				JoinTime:        membersGJson.Get(qq).Get("jt").Int(),
				LastSentTime:    membersGJson.Get(qq).Get("lst").Int(),
				Level:           "unknown",
				Role:            role,
				Unfriendly:      false,
				Title:           "unknown",
				TitleExpireTime: 0,
				CardChangeable:  true,
			}
			group.GroupMembers = append(group.GroupMembers, member)
		}
		c.M.Lock()
		c.Group = append(c.Group, group)
		c.M.Unlock()
		return group
	}
	listC := XQApiGroupMemberListC(bot, groupID)
	if listC != "" {
		group.GroupInfo = &GroupInfo{
			GroupID:        groupID,
			GroupName:      name,
			MemberCount:    0,
			MaxMemberCount: 0,
		}
		qqMap := gjson.Parse(listC).Get("list").Map()
		for _, value := range qqMap {
			qq := value.Get("QQ").Str
			member := &GroupMember{
				GroupID:         groupID,
				UserID:          str2Int(qq),
				Nickname:        "unknown",
				Card:            "unknown",
				Sex:             "unknown",
				Age:             0,
				Area:            "unknown",
				JoinTime:        0,
				LastSentTime:    0,
				Level:           "unknown",
				Role:            "unknown",
				Unfriendly:      false,
				Title:           "unknown",
				TitleExpireTime: 0,
				CardChangeable:  true,
			}
			group.GroupMembers = append(group.GroupMembers, member)
		}
		c.M.Lock()
		c.Group = append(c.Group, group)
		c.M.Unlock()
		return group
	}
	return nil
}

// GetCacheGroupMember 获取群成员信息，不存在则请求
func (c *CacheGroupsData) GetCacheGroupMember(bot, groupID, userID int64, cache bool) (member *GroupMember) {
	members := c.GetCacheGroup(bot, groupID, cache).GroupMembers
	if members == nil {
		return nil
	}
	for i := range members {
		if userID == members[i].UserID {
			member = members[i]
			break
		}
	}
	if member == nil {
		return nil
	}

	if member.Nickname == "" {
		c.M.Lock()
		member.Nickname = XQApiGetNick(bot, userID)
		c.M.Unlock()
	}
	if member.Age == 0 {
		c.M.Lock()
		member.Age = XQApiGetAge(bot, userID)
		c.M.Unlock()
	}
	if member.Sex == "" {
		c.M.Lock()
		member.Sex = XQApiGetGender(bot, userID)
		c.M.Unlock()
	}
	return member
}

func (m *GroupMember) GetNick() string {
	if m == nil {
		return "unknown"
	}
	return m.Nickname
}

func (m *GroupMember) GetAge() int64 {
	if m == nil {
		return 0
	}
	return m.Age
}

func (m *GroupMember) GetSex() string {
	if m == nil {
		return "unknown"
	}
	return m.Sex
}

func (m *GroupMember) GetRole() string {
	if m == nil {
		return "member"
	}
	return m.Role
}
