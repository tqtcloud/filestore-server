create database fileserver default character set utf8;

create table `tbl_file` (
                            `id` int(11) not null auto_increment ,
                            `file_sha1` char(40) not null default '' comment '文件hash',
                            `file_name` varchar(256) not null default '' comment '文件名',
                            `file_size` bigint(20) default '0' comment '文件大小',
                            `file_addr` varchar(1024) not null default '' comment '文件存储位置',
                            `create_at` datetime DEFAULT CURRENT_TIMESTAMP  comment '创建日期',
                            `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  comment '更新时间',
                            `delete_at` datetime comment '更新时间',
                            `status` int(11) not null default '0' comment '状态(可用/禁用/已删除)',
                            `ext1` int(11) default '0' comment '备用字段1',
                            `ext2` text comment '备用字段2' ,
                            primary key (`id`) ,
                            KEY `idx_status` (`status`),
                            unique key `idx_file_hash`(`file_sha1`)
) engine=innodb default charset=utf8;


show variables like 'sql_mode';
set global sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';