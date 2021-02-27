
create table if not exists videos (
	video_id int unsigned primary key auto_increment,
    path varchar(255) not null,
    post_id int unsigned not null,
    foreign key fk_videos_posts (post_id)
		references posts (post_id)
        on update cascade
        on delete cascade
);