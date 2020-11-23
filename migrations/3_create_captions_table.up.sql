create table captions
(
    id               uuid          not null,
    video_id         uuid          not null
        constraint captions_videos_id_fk
            references videos
            on delete cascade,
    text             varchar       not null,
    start            int default 0 not null,
    "end"            int default 0 not null,
    previous_caption uuid,
    next_caption     uuid
);

create unique index captions_id_uindex
    on captions (id);

create index captions_video_id_index
    on captions (video_id);

alter table captions
    add constraint captions_pk
        primary key (id);
