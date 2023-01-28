# 最大长度
#   用户名(nickname)和用户登陆名(username):  100
#   密码(password): 200
#   url:
#       视频播放地址(play_url): 200
#       封面地址(cover_url): 200

# 索引
# follow_list:
#       (from_user_uuid, to_user_uuid):用户关注的所有用户列表, 查询是否关注的时候使用
#       (to_user_uuid,from_user_uuid): 所有关注登录用户的粉丝列表
# user_favourite_video:
#     (uuid, video_id, is_like): 查询喜欢列表的时候使用
# comment_info:
#       (video_id): 查看视频的所有评论.

create database if not exists dy;
CREATE TABLE if not exists dy.user_info
(
    uuid           int8         NOT NULL COMMENT '用户唯一标识',
    username       varchar(100) NOT NULL comment '用户登陆名',
    password       varchar(200) NOT NULL COMMENT '用户登陆密码',
    follow_count   int DEFAULT 0 COMMENT '用户关注数量',
    follower_count int DEFAULT 0 COMMENT '用户粉丝数量',
    nickname       varchar(100) NOT NULL COMMENT '用户名',
    created_at     date         NOT NULL,
    updated_at     date         NOT NULL,
    deleted_at     date,
    PRIMARY KEY (uuid),
    UNIQUE KEY (username),
    check ( follow_count >= 0 ),
    check ( follower_count >= 0 )
) ENGINE = InnoDB
  DEFAULT charset = utf8mb4;

# 关注列表
CREATE TABLE if not exists dy.follow_list
(
    from_user_uuid int8 NOT NULL COMMENT '发起关注请求的人',
    to_user_uuid   int8 NOT NULL COMMENT '被关注的人',
    created_at     date NOT NULL,
    updated_at     date NOT NULL,
    deleted_at     date,
    UNIQUE KEY (from_user_uuid, to_user_uuid) COMMENT '查询是否关注的时候使用以及用户关注的所有用户列表 (from,to)',
    UNIQUE KEY (to_user_uuid, from_user_uuid) COMMENT '所有关注登录用户的粉丝列表'
) ENGINE = InnoDB
  DEFAULT charset = utf8mb4;

# 用户喜欢的视频表
CREATE TABLE if not exists dy.user_favourite_video
(
    uuid       int8 NOT NULL COMMENT '用户唯一标识',
    video_id   int  NOT NULL,
    is_like    tinyint DEFAULT 0 COMMENT '喜欢:0, 不喜欢:1',
    created_at date NOT NULL,
    updated_at date NOT NULL,
    deleted_at date,
    KEY (uuid, video_id, is_like) COMMENT '查询喜欢列表的时候使用 (from,to,0/1)'
) ENGINE = InnoDB
  DEFAULT charset = utf8mb4;


# 视频信息表
CREATE TABLE if not exists dy.video_info
(
    video_id        int          NOT NULL,
    play_url        varchar(200) NOT NULL COMMENT '视频地址',
    uuid            int8         NOT NULL COMMENT '用户唯一标识',
    cover_url       varchar(200) NOT NULL COMMENT '封面地址',
    title           varchar(100) NOT NULL COMMENT '视频标题',
    favourite_count int          NOT NULL COMMENT '视频点赞数',
    comment_count   int          NOT NULL COMMENT '视频评论数',
    created_at      date         NOT NULL,
    updated_at      date         NOT NULL,
    deleted_at      date,
    PRIMARY KEY (video_id)
) ENGINE = InnoDB
  DEFAULT charset = utf8mb4;

# 评论信息表
CREATE TABLE if not exists dy.comment_info
(
    comment_id int  NOT NULL COMMENT '评论id,删除的时候使用',
    uuid       int8 NOT NULL COMMENT '用户唯一标识',
    video_id   int  NOT NULL COMMENT '视频id',
    content    text NOT NULL COMMENT '评论内容',
    created_at date NOT NULL,
    updated_at date NOT NULL,
    deleted_at date,
    PRIMARY KEY (`comment_id`),
    KEY (video_id) COMMENT '查询视频评论列表使用'
) ENGINE = InnoDB
  DEFAULT charset = utf8mb4;
