create table videos
(
    id          uuid    not null,
    video_id    varchar not null,
    caption     varchar,
    has_caption bool    not null,
    processed   bool    not null
);

create unique index videos_id_uindex
    on videos (id);

create unique index videos_video_id_uindex
    on videos (video_id);

alter table videos
    add constraint videos_pk
        primary key (id);
