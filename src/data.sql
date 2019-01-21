/*创建数据库*/
CREATE DATABASE IF NOT EXISTS `bg_cms`;
USE `bg_cms`;

/*用户表*/
CREATE TABLE IF NOT EXISTS `bc_users`(
    `user_id`  BIGINT(20)UNSIGNED AUTO_INCREMENT,
    `user_login` VARCHAR(60) NOT NULL COMMENT '用户名',
    `user_pass` VARCHAR(255) NOT NULL COMMENT '密码',
    `user_nicename` VARCHAR(50) DEFAULT '山猫' COMMENT '昵称',
    `user_email` VARCHAR(50)  DEFAULT '' COMMENT '邮箱',
    `user_registered` DATETIME DEFAULT now() COMMENT '注册时间',
    `user_activation_key` VARCHAR(60)  DEFAULT '0000' COMMENT '激活码',
    `user_status`   INT DEFAULT 0 COMMENT '用户状态(0正常，1不正常)',
    PRIMARY KEY (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 用户元素据表 */
CREATE TABLE IF NOT EXISTS `bc_usermeta`(
    `meta_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `user_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键用户ID',
    `meta_key` VARCHAR(255),
    `meta_value` LONGTEXT,
    PRIMARY KEY (`meta_id`),
    foreign key(`user_id`) references bc_users(`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*
bc_terms -----1:N-------> bc_term_taxonomy<------N:1----bc_term_relationships<-----N:1----- bc_posts/bc_links
bc_terms记录了每个分类的名字以及基本信息，如分为“JAVA开发”、“HTML5”等，
这里的分类指广义上的分类，所以每个TAG也是一个“分类”。bc_term_taxonomy记录了每个分类所归属的分类方法，
如“JAVA开发”、“HTML5”是文章分类（category），放置友情链接的“合作商”、“我的好友”分类属于友情链接分类（link_category）。
bc_term_relationships记录了每个文章（或链接）所对应的分类方法。
*/
/* 分类表 表存储目录，标签，链接目录和自定义分类的所有单个分类项*/
CREATE TABLE IF NOT EXISTS `bc_terms`(
    `term_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `term_name` varchar(200) NOT NULL  COMMENT '分类名',
    `term_group` INT DEFAULT 0 COMMENT '分类组',
    PRIMARY KEY (`term_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 分类法表 存储分类和标签的描述信息、父子关系 存储每个目录、标签所对应的分类
表存储更多关于分类项的数据以及他们属于的分类的 信息  可以理解为【分类项信息表】 对分类项进行拓展
分类法：
    一个分类法是一个目录化或分类化事物的系统，通常以分级的方式进行。最著名的分类法是 Linnean Taxonomy ,用来对所有活的的事物分类
*/
CREATE TABLE IF NOT EXISTS `bc_term_taxonomy`(
    `term_taxonomy_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `term_id` BIGINT(20)UNSIGNED NOT NULL COMMENT '外键分类ID',
    `taxonomy` VARCHAR(32) NOT NULL COMMENT '分类法名',
    `description` VARCHAR(255) COMMENT '分类法描述',
    `term_parent_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '父分类ID',
    PRIMARY KEY (`term_taxonomy_id`),
    foreign key(`term_id`) references bc_terms(`term_id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 存储文章和分类、标签的相互对应关系*/
CREATE TABLE IF NOT EXISTS `bc_term_relationships`(
    `object_id` BIGINT(20)UNSIGNED DEFAULT 0  COMMENT '文章ID/链接ID',
    `term_taxonomy_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键分类法ID',
    foreign key(`term_taxonomy_id`) references bc_term_taxonomy(`term_taxonomy_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 文章表 */
CREATE TABLE IF NOT EXISTS `bc_posts`(
    `post_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `post_author`  BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '作者',
    `post_date`  DATETIME DEFAULT now() COMMENT '时间',
    `post_content` LONGTEXT COMMENT '内容',
    `post_title` VARCHAR(255) COMMENT '标题',
    `post_excerpt` TEXT COMMENT '描述',
    `post_pre_img_url` VARCHAR(255) COMMENT '文章预览图',
    `post_status` INT DEFAULT 0 COMMENT '文章状态(0发表,1不发表)',
    `comment_status` INT DEFAULT 0 COMMENT '评论状态(0允许,1不允许)',
    `post_modified` DATETIME ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    `comment_count` INT DEFAULT 0 COMMENT '评论数',
    PRIMARY KEY (`post_id`),
    foreign key(`post_author`) references bc_users(`user_id`) ON DELETE SET NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 文章元数据表 */
CREATE TABLE IF NOT EXISTS `bc_postmeta`(
    `meta_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `post_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键文章ID',
    `meta_key` VARCHAR(255),
    `meta_value` LONGTEXT,
    PRIMARY KEY (`meta_id`),
    foreign key(`post_id`) references bc_posts(`post_id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


/* 评论表 */
CREATE TABLE IF NOT  EXISTS `bc_comments`(
    `comment_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `comment_post_id` BIGINT(20)UNSIGNED COMMENT '外键文章ID',
    `comment_author` VARCHAR(255) COMMENT '评论作者',
    `comment_author_email` VARCHAR(100) COMMENT '评论作者邮箱',
    `comment_author_IP` VARCHAR(100) COMMENT '评论作者IP',
    `comment_date` DATETIME DEFAULT now() COMMENT '评论时间',
    `comment_content` TEXT COMMENT '评论内容',
    `comment_approved` INT DEFAULT 1 COMMENT '是否被允许(0是，1否)',
    `comment_parent` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键父评论ID',
    `user_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键用户ID',
    PRIMARY KEY (`comment_id`),
    foreign key(`comment_post_id`) references bc_posts(`post_id`) ON DELETE CASCADE,
    foreign key(`user_id`) references bc_users(`user_id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 评论元数据表 */
CREATE TABLE IF NOT  EXISTS `bc_commentmeta`(
    `meta_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `comment_id` BIGINT(20)UNSIGNED DEFAULT 0 COMMENT '外键评论ID',
    `meta_key` VARCHAR(255),
    `meta_value` LONGTEXT,
    PRIMARY KEY (`meta_id`),
    foreign key(`comment_id`) references bc_comments(`comment_id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 连接表 */
CREATE TABLE IF NOT  EXISTS `bc_links`(
    `link_id` BIGINT(20)UNSIGNED  AUTO_INCREMENT,
    `link_url` VARCHAR(255) COMMENT '链接URL',
    `link_name` VARCHAR(255) COMMENT '链接名称',
    `link_image_url` VARCHAR(255) COMMENT '链接图片URL',
    `link_description` VARCHAR(255) COMMENT '链接描述',
    `link_visible` INT DEFAULT 0 COMMENT '是否可见(0可见，1不可见)',
    `link_owner` BIGINT(20)UNSIGNED DEFAULT 1 COMMENT '添加者ID',
    `link_rating` INT DEFAULT 0 COMMENT '链接评分等级',
    `link_updated` DATETIME ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP  COMMENT '链接更新时间',
    PRIMARY KEY (`link_id`),
    foreign key(`link_owner`) references bc_users(`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/* 创建超级管理员账户 */
INSERT INTO bc_users (user_login,user_pass) VALUES ('admin','123456');