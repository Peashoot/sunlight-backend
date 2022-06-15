package utils

const (
	UserLoginLockPrefix        = "user_login_lock_key_"     // redis用户登录分布式锁前缀
	UserRegLockPrefix          = "user_reg_lock_key_"       // redis用户注册分布式锁前缀
	UserCagKeyLockPrefix       = "user_cag_lock_key_"       // redis用户修改分布式锁前缀
	UserDelKeyLockPrefix       = "user_del_lock_key_"       // redis用户删除分布式锁前缀
	LabelAddKeyLockPrefix      = "label_add_lock_key_"      // redis标签新增分布式锁前缀
	UnvalidUserCachePrefix     = "unvalid_user_cache_key_"  // redis无效账号缓存前缀
	UserRedisCachePrefix       = "user_entity_cache_key_"   // redis用户信息缓存前缀
	JwtRedisCachePrefix        = "user_auth_cache_key_"     // redis鉴权信息缓存前缀
	UserLabelCachePrefix       = "user_label_cache_key_"    // redis用户标签缓存前缀
	UserGroupCachePrefix       = "user_group_cache_key_"    // redis用户分组缓存前缀
	GroupMembershipCachePrefix = "group_memship_cache_key_" // redis群组关系缓存前缀
	FileCategoryCachePrefix    = "file_category_cache_key_" // redis文件分类缓存前缀
	EnumDictionaryCachePrefix  = "enum_dict_cache_key_"     // redis枚举字典缓存前缀
)
