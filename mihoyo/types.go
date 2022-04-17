package mihoyo

type TaskListResp struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		States []struct {
			MissionID     int    `json:"mission_id"`
			Process       int    `json:"process"`
			HappenedTimes int    `json:"happened_times"`
			IsGetAward    bool   `json:"is_get_award"`
			MissionKey    string `json:"mission_key"`
		} `json:"states"`
		AlreadyReceivedPoints int  `json:"already_received_points"`
		TotalPoints           int  `json:"total_points"`
		TodayTotalPoints      int  `json:"today_total_points"`
		IsUnclaimed           bool `json:"is_unclaimed"`
		CanGetPoints          int  `json:"can_get_points"`
	} `json:"data"`
}

type PostListResp struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Post struct {
				GameID     int           `json:"game_id"`
				PostID     string        `json:"post_id"`
				FForumID   int           `json:"f_forum_id"`
				UID        string        `json:"uid"`
				Subject    string        `json:"subject"`
				Content    string        `json:"content"`
				Cover      string        `json:"cover"`
				ViewType   int           `json:"view_type"`
				CreatedAt  int           `json:"created_at"`
				Images     []interface{} `json:"images"`
				PostStatus struct {
					IsTop      bool `json:"is_top"`
					IsGood     bool `json:"is_good"`
					IsOfficial bool `json:"is_official"`
				} `json:"post_status"`
				TopicIds               []int         `json:"topic_ids"`
				ViewStatus             int           `json:"view_status"`
				MaxFloor               int           `json:"max_floor"`
				IsOriginal             int           `json:"is_original"`
				RepublishAuthorization int           `json:"republish_authorization"`
				ReplyTime              string        `json:"reply_time"`
				IsDeleted              int           `json:"is_deleted"`
				IsInteractive          bool          `json:"is_interactive"`
				StructuredContent      string        `json:"structured_content"`
				StructuredContentRows  []interface{} `json:"structured_content_rows"`
				ReviewID               int           `json:"review_id"`
				IsProfit               bool          `json:"is_profit"`
				IsInProfit             bool          `json:"is_in_profit"`
				UpdatedAt              int           `json:"updated_at"`
				DeletedAt              int           `json:"deleted_at"`
				PrePubStatus           int           `json:"pre_pub_status"`
			} `json:"post"`
			Forum struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Icon   string `json:"icon"`
				GameID int    `json:"game_id"`
			} `json:"forum"`
			Topics []struct {
				ID            int    `json:"id"`
				Name          string `json:"name"`
				Cover         string `json:"cover"`
				IsTop         bool   `json:"is_top"`
				IsGood        bool   `json:"is_good"`
				IsInteractive bool   `json:"is_interactive"`
				GameID        int    `json:"game_id"`
				ContentType   int    `json:"content_type"`
			} `json:"topics"`
			User struct {
				UID           string `json:"uid"`
				Nickname      string `json:"nickname"`
				Introduce     string `json:"introduce"`
				Avatar        string `json:"avatar"`
				Gender        int    `json:"gender"`
				Certification struct {
					Type  int    `json:"type"`
					Label string `json:"label"`
				} `json:"certification"`
				LevelExp struct {
					Level int `json:"level"`
					Exp   int `json:"exp"`
				} `json:"level_exp"`
				IsFollowing bool   `json:"is_following"`
				IsFollowed  bool   `json:"is_followed"`
				AvatarURL   string `json:"avatar_url"`
				Pendant     string `json:"pendant"`
			} `json:"user"`
			SelfOperation struct {
				Attitude    int  `json:"attitude"`
				IsCollected bool `json:"is_collected"`
			} `json:"self_operation"`
			Stat struct {
				ViewNum     int `json:"view_num"`
				ReplyNum    int `json:"reply_num"`
				LikeNum     int `json:"like_num"`
				BookmarkNum int `json:"bookmark_num"`
				ForwardNum  int `json:"forward_num"`
			} `json:"stat"`
			HelpSys struct {
				TopUp     interface{}   `json:"top_up"`
				TopN      []interface{} `json:"top_n"`
				AnswerNum int           `json:"answer_num"`
			} `json:"help_sys"`
			Cover            interface{}   `json:"cover"`
			ImageList        []interface{} `json:"image_list"`
			IsOfficialMaster bool          `json:"is_official_master"`
			IsUserMaster     bool          `json:"is_user_master"`
			HotReplyExist    bool          `json:"hot_reply_exist"`
			VoteCount        int           `json:"vote_count"`
			LastModifyTime   int           `json:"last_modify_time"`
			RecommendType    string        `json:"recommend_type"`
			Collection       interface{}   `json:"collection"`
			VodList          []interface{} `json:"vod_list"`
			IsBlockOn        bool          `json:"is_block_on"`
			ForumRankInfo    interface{}   `json:"forum_rank_info"`
			LinkCardList     []interface{} `json:"link_card_list"`
		} `json:"list"`
		LastID   string `json:"last_id"`
		IsLast   bool   `json:"is_last"`
		IsOrigin bool   `json:"is_origin"`
	} `json:"data"`
}

type GenshinAccountsResp struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			GameBiz    string `json:"game_biz"`
			Region     string `json:"region"`
			GameUID    string `json:"game_uid"`
			Nickname   string `json:"nickname"`
			Level      int    `json:"level"`
			IsChosen   bool   `json:"is_chosen"`
			RegionName string `json:"region_name"`
			IsOfficial bool   `json:"is_official"`
		} `json:"list"`
	} `json:"data"`
}

type GenshinSignInfoResp struct {
	Retcode int             `json:"retcode"`
	Message string          `json:"message"`
	Data    GenshinSignInfo `json:"data"`
}

type GenshinSignInfo struct {
	TotalSignDay  int    `json:"total_sign_day"`
	Today         string `json:"today"`
	IsSign        bool   `json:"is_sign"`
	FirstBind     bool   `json:"first_bind"`
	IsSub         bool   `json:"is_sub"`
	MonthFirst    bool   `json:"month_first"`
	SignCntMissed int    `json:"sign_cnt_missed"`
}

type GenshinSignPostData struct {
	ActID  string `json:"act_id"`
	Region string `json:"region"`
	UID    string `json:"uid"`
}

type HomuShopGoodListResp struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		List  []*HomuShopGood `json:"list"`
		Total int             `json:"total"`
		Games []struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"games"`
	} `json:"data"`
}

type HomuShopGood struct {
	AppID              int    `json:"app_id"`
	GoodsID            string `json:"goods_id"`
	GoodsName          string `json:"goods_name"`
	Type               int    `json:"type"`
	Price              int    `json:"price"`
	PointSn            string `json:"point_sn"`
	Icon               string `json:"icon"`
	Unlimit            bool   `json:"unlimit"`
	Total              int    `json:"total"`
	AccountCycleType   string `json:"account_cycle_type"`
	AccountCycleLimit  int    `json:"account_cycle_limit"`
	AccountExchangeNum int    `json:"account_exchange_num"`
	RoleCycleType      string `json:"role_cycle_type"`
	RoleCycleLimit     int    `json:"role_cycle_limit"`
	RoleExchangeNum    int    `json:"role_exchange_num"`
	Start              string `json:"start"`
	End                string `json:"end"`
	Status             string `json:"status"`
	NextTime           int    `json:"next_time"`
	NextNum            int    `json:"next_num"`
	NowTime            int    `json:"now_time"`
}
