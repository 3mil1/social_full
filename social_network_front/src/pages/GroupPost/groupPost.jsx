import { Post } from "../../components/posts/Post";
import { useParams } from "react-router-dom";
import { Container } from "@mui/material";
import { useEffect, useState } from "react";
import GroupService from "../../utilities/group_service";
import Create_comment from "../../components/groups/buttons_forms/Create_comment";

import "../../components/posts/post.scss";

export default function GroupPost() {
  const group_service = GroupService();
  let { groupId, postId } = useParams();
  const group_id=  groupId  ? groupId : "";
  const post_id = postId ? postId : "";
  const [update,setUpdate] = useState(false)
  const [post,setPost] = useState();
  const [comments, setComments] = useState();


  const handleComment = () => {
    console.log("HAndeling submit");
    setUpdate(!update)
  }

  useEffect(()=>{
    // if (post) return;
    group_service.getSpecificGroupPost(group_id,post_id).then(data => { 
    const p = data.Post;
    const c = data.Comments || [];
    const post = {
          id: p.post_id,
          user_id: p.user_id,
          user_firstname: p.user_firstname,
          user_lastname: p.User_lastname,
          title: p.subject,
          content: p.content,
          image: p.image,
          created_at: p.created_at,
        };
    const comments= [];
    c.forEach((v) => {
      const com = {
        id: v.post_id,
        user_id: v.user_id,
        user_firstname: v.user_firstname,
        user_lastname: v.user_lastname,
        title: v.subject,
        content: v.content,
        image: v.image,
        created_at: v.created_at,
      };
      comments.push(com);
    });
    setPost(post)
    setComments(comments);
    })
  },[update])
 
  if (!post) {
    return <div>...Loading</div>;
  }

  return (
    <Container>
      <div className="post">
        <Post post={post} toShow={false} ></Post>
        <Create_comment className="create_comment_btn" group_id={group_id} post_id={post_id}  handleComment={handleComment}/> 
      </div>
      <div className="post_list" style={{ maxWidth: 600 }}>
        {comments && comments.map((c) => (
          <Post key={c.id} post={c} toShow={false} />
        ))}
      </div>
    </Container>
  );
}
