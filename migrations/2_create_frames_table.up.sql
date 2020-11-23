create table frames
(
    id             uuid          not null,
    video_id       uuid          not null
        constraint frames_videos_id_fk
            references videos
            on delete cascade,
    file_name      varchar,
    time           int default 0 not null,
    image          bytea,
    previous_frame varchar
);

create unique index frames_id_uindex
    on frames (id);

create index frames_video_id_time_index
    on frames (video_id, time);


create index frames_previous_frame_index
    on frames (previous_frame);

alter table frames
    add constraint frames_pk
        primary key (id);
