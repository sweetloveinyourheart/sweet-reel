export type ChannelResponse = {
  id: string;
  owner_id: string;
  name: string;
  handle: string;
  description: string;
  banner_url: string;
  subscriber_count: number;
  total_views: number;
  total_videos: number;
  created_at: string;
  updated_at: string;
  owner: OwnerResponse;
};

type OwnerResponse = {
  id: string;
  email: string;
  name: string;
  picture: string;
  created_at: string;
  updated_at: string;
};