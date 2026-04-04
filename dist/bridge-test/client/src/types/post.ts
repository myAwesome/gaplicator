export interface Post {
  id: number;
  title?: string;
}

export type CreatePostInput = {
  title: string;
};
