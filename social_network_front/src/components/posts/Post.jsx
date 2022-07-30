import "./post.scss";
import { Button } from "@mui/material";
import { Link } from "react-router-dom";
import { openModal } from "../../store/postSlice";
import { useDispatch } from "react-redux";
// import * as React from "@types/react";

export const Post = ({ post, toShow }) => {
  const dispatch = useDispatch();

  const handleClick = () => {
    dispatch(openModal());
  };

  return (
    <div className="post">
      <div className="post_header">
        {/*<img src={require("../../assets/Images/ano.jpg")} alt="ano_pic" />*/}
        {post.title !== "" ? (
          <>
            <div className="information">
              <div className="name title">
                {toShow ? (
                  <div>{post.title}</div>
                ) : (
                  <Link to={`/post/${post.id}`}>{post.title}</Link>
                )}
              </div>
              <div className="user-name">
                {post.user_firstname} {post.user_lastname}
              </div>
            </div>
            <div className="date">{post.created_at}</div>
          </>
        ) : (
          <>
            <div className="information">
              <div className="name">
                {" "}
                {post.user_firstname} {post.user_lastname}
              </div>
              <div className="date">{post.created_at}</div>
            </div>
          </>
        )}
      </div>
      <div className="post_content">
        <div>{post.content}</div>
        {post.image && <img src={`${post.image}`} />}
      </div>
      {toShow && (
        <Button
          onClick={(e) => {
            e.preventDefault();
            handleClick();
          }}
        >
          Comment
        </Button>
      )}
    </div>
  );
};

export default Post;
