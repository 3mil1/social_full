import { Post } from "../../components/posts/Post";
import { useParams } from "react-router-dom";
import { Container } from "@mui/material";
import React, { useEffect, useState } from "react";
import { PostInterface } from "../../components/posts/PostList";
import "../../components/posts/post.scss";
import postService from "../../utilities/post-service";
import { NewPost } from "../../components/posts/newPost";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../../store/store";
import { loadComments } from "../../store/postSlice";
import { parseDate } from "../../helpers/parseDate";

export default function OnePost() {
  let { id } = useParams();
  const postId = id ? id : "";
  const [post, setPost] = React.useState<PostInterface>();
  // const [comments, setComments] = React.useState<PostInterface[]>([]);
  const isOpen = useSelector((state: RootState) => state.post.isOpen);

  const dispatch = useDispatch();
  const comments: PostInterface[] = useSelector(
    (state: RootState) => state.post.comments
  );

  useEffect(() => {
    if (post) return;
    const getPostComments = async (id: string) => {
      try {
        // @ts-ignore
        const data = await postService.showPost(id);
        const p = data.Post;
        console.log("P", p);
        const c = data.Comments || [];
        console.log("C", c);
        const date = parseDate(p.created_at);
        const post: PostInterface = {
          id: p.id,
          user_id: p.user_id,
          user_firstname: p.user_firstname,
          user_lastname: p.user_lastname,
          title: p.subject,
          content: p.content,
          image: p.image,
          privacy: p.privacy,
          created_at: date,
        };
        console.log("POST", post);
        const comments: PostInterface[] = [];
        c.forEach((v: any) => {
          const com = {
            id: v.id,
            user_id: v.user_id,
            user_firstname: v.user_firstname,
            user_lastname: v.user_lastname,
            title: v.subject,
            content: v.content,
            image: v.image,
            created_at: date,
            privacy: v.privacy,
          };
          comments.push(com);
        });
        setPost(post);
        // setComments(comments);
        dispatch(loadComments(comments));
      } catch (e) {
        console.error(e);
      }
    };
    getPostComments(postId);
  });

  if (!post) {
    return <div>...Loading</div>;
  }

  return (
    <Container>
      <Post post={post} toShow={true}></Post>
      <div className="post_list" style={{ maxWidth: 600 }}>
        {comments.map((c) => (
          <Post key={c.id} post={c} toShow={false} />
        ))}
      </div>
      {isOpen ? <NewPost parentPrivacy={post.privacy} /> : null}
    </Container>
  );
}
