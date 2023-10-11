/// 用户对象
#[derive(Debug)]
pub struct User {
    /// 用户 id
    pub id: String,

    /// 用户名
    pub username: String,

    /// 用户头像地址
    pub avatar: String,

    /// 是否是机器人
    pub bot: bool,

    /// 特殊关联应用的 openid，
    /// 需要特殊申请并配置后才会返回。
    /// 如需申请，请联系平台运营人员。
    pub union_openid: String,

    /// 机器人关联的互联应用的用户信息，
    /// 与union_openid关联的应用是同一个。
    /// 如需申请，请联系平台运营人员。
    pub union_user_account: String,
}
