export interface Media {
    id: number;
    user_id?: number;
    slug: string;
    title: string;
    type: 'image' | 'video' | 'document' | 'audio' | 'other';
    url: string;
    file_name: string;
    file_size: number;
    mime_type: string;
    visibility: 'public' | 'private' | 'shared';
    created_at: string;
    updated_at: string;
}

export interface Movie {
    id: number;
    user_id?: number;
    media_id: number;
    title: string;
    description?: string;
    release_date?: string;
    duration?: number;
    created_at: string;
    updated_at: string;
}

export interface Series {
    id: number;
    user_id?: number;
    title: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface Season {
    id: number;
    series_id: number;
    season_number: number;
    title?: string;
    created_at: string;
    updated_at: string;
}

export interface Episode {
    id: number;
    season_id: number;
    media_id: number;
    episode_number: number;
    title?: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface MusicAlbum {
    id: number;
    user_id?: number;
    title: string;
    artist?: string;
    cover_media_id?: number;
    release_date?: string;
    created_at: string;
    updated_at: string;
}

export interface MusicTrack {
    id: number;
    album_id: number;
    media_id: number;
    title: string;
    track_number?: number;
    duration?: number;
    created_at: string;
    updated_at: string;
}
