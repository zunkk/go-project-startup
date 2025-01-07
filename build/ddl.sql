-- 用户信息
create table if not exists "user"
(
    "id"                               bigint       not null
        constraint user_pk
            primary key,
    "create_time"                      timestamptz  not null,
    "update_time"                      timestamp    not null,
    "delete_time"                      timestamp    not null,
    -- 删除状态(使用软删除)
    -- 0:active 1:deleted
    "del_state"                        bigint       not null default 0,
    "version"                          bigint       not null default 0,

    -- 用户名
    "username"                         varchar(255) not null default '',
    -- 用户信息
    "info"                             varchar(255) not null default '',
    -- 角色
    "role"                             varchar(20)  not null default '',
    -- 最后一次使用的钱包发送交易的配置id
    last_used_wallet_send_tx_config_id bigint       not null default 0
);

-- 用户认证信息
-- 关联关系: user - user_auth ：1 - n
create table if not exists "user_auth"
(
    "id"              bigint       not null
        constraint user_auth_pk
            primary key,
    "create_time"     timestamptz  not null,
    "update_time"     timestamp    not null,
    "delete_time"     timestamp    not null,
    -- 0:active 1:deleted
    "del_state"       bigint       not null default 0,
    "version"         bigint       not null default 0,

    -- 关联的用户id
    "user_id"         bigint       not null default 0,
    -- 认证类型
    -- tg email
    "auth_type"       varchar(20)  not null default '',
    -- 认证渠道的id
    -- tg: tg的用户id
    -- email: 邮箱地址
    "auth_id"         varchar(64)  not null default '',
    -- 认证渠道的token
    -- tg: 无(因为消息走tg，tg已经做完这一步认证了)
    -- email: 密码
    "auth_token"      varchar(255) not null default '',
    -- 上一次登录时间
    "last_login_time" timestamp    not null
);

create index if not exists user_auth_type_index on "user_auth" ("auth_type", "auth_id");
create index if not exists user_auth_user_id_index on "user_auth" ("user_id", "auth_type");