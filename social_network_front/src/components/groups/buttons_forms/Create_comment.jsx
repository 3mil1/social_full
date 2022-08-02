import { Button, Input } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";
import GroupService from "../../../utilities/group_service";
import * as helper from "../../../helpers/HelperFuncs";
import "./group_buttons.scss";

const Create_comment = ({ group_id, post_id, handleComment }) => {
  const group_service = GroupService();
  const [isOpen, setIsOpen] = useState(false);
  const [img, setImg] = useState(null);
  const [content, setContent] = useState(null);

  const data = {
    group_id: Number(group_id),
    parent_id: Number(post_id),
    content: content,
    image: img,
  };

  const convertImg = (image) => {
    if (image.length !== 0) {
      if (helper.checkImage(image)) {
        const resp =  helper.getBase64(image[0]).then((base64) => base64);
        return resp;
      }
    }
  };

  const clearInput = (e) => {
    e.target.value = "";
    document.getElementById(e.target.id).classList.remove("error");
  };

  const handleSubmit = () => {
    if (data == null) return;
    if (helper.handleInputs("content", data.content)) {
      group_service.makeCommentToPost(data);
      setIsOpen(!isOpen);
    }
    handleComment();
  };
  
  
  return (
    <div className="comment_wrapper">
      {
        <Button
        sx={{ marginLeft: 3, marginBottom: 1, fontSize: "18px" }}
        onClick={() => {
            setIsOpen(!isOpen)
            setImg(null)
            setContent(null)
          }
          }
        >
          Comment
        </Button>
      }

      {isOpen && (
        <form id="commentForm">
          <Button
            variant="contained"
            className="back_btn"
            onClick={() => setIsOpen(false)}
            >
            <CloseIcon />
          </Button>
          <div className="input">
            <label htmlFor="content">Content* : </label>
            <Input
              type="text"
              id="content"
              onClick={(e) => clearInput(e)}
              onChange={(e)=>{setContent(e.target.value)}}
              ></Input>
          </div>
          <div className="input">
            <label className="image_btn" htmlFor="image">
              {!img ? "PICK IMAGE" : "IMAGE ADDED"}
            </label>
            <input
              type="file"
              id="image"
              name="image"
              onChange={() => {
                convertImg(document.getElementById("image").files).then(res => {
                  setImg(res)
                })
            }}
            />
          </div>

          <Button
            sx={{ fontSize: "16px" }}
            type={"submit"}
            onClick={(e) => {
              e.preventDefault();
              if (data.content) setIsOpen(false);
              handleSubmit();
            }}
            >
            COMMENT
          </Button>
        </form>
      )}
    </div>
  );
};

export default Create_comment;
