import {
  Box,
  Button,
  FormControl,
  Input,
  Modal,
  Radio,
  RadioGroup,
  TextField,
} from "@mui/material";
import * as React from "react";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { openModal, updateComments, updatePosts } from "../../store/postSlice";
import "./newPost.scss";
import CloseIcon from "@mui/icons-material/Close";
import FormControlLabel from "@mui/material/FormControlLabel";
import { RootState } from "../../store/store";
import { useForm } from "react-hook-form";
import { checkImage, getBase64 } from "../../helpers/checkImage";
import postService from "../../utilities/post-service";
import { useParams } from "react-router-dom";
import { setAlert } from "../../store/alertSlice";
import getTokenId from "../../helpers/tokenId";
import { UserList } from "./checkBox";

const ariaLabel = { "aria-label": "description" };

export interface NewPostForm {
  title: string;
  content: string;
  image: FileList;
  imgString: string;
  privacy: Privacy;
  parent_id: number;
  userList: string[] | null;
}

export interface Follower {
  firstName: string;
  lastName: string;
  id: string;
}

interface FollowerFromState {
  first_name: string;
  last_name: string;
  user_id: string;
}

enum Privacy {
  Public = 1,
  Private,
  StrictlyPrivate,
}

export function NewPost(props: { parentPrivacy: number }) {
  const myFollowers = useSelector(
    (state: RootState) => state.followers.stalkers
  );
  const { handleSubmit, register } = useForm<NewPostForm>();
  const [errors, setErrors] = useState<string[]>([]);
  const open = useSelector((state: RootState) => state.post.isOpen);
  const dispatch = useDispatch();
  const handleClose = () => {
    dispatch(openModal());
  };
  const [value, setValue] = React.useState(Privacy.Public);
  const [followers, setFollowers] = React.useState<boolean>(false);
  const [listOfFollowers, setListFollowers] = React.useState<Follower[]>([]);

  let { id } = useParams();
  const param: number = id ? +id : 0;
  
  useEffect(() => {
    if (listOfFollowers.length !== 0) {
      return;
    }
    let l: Follower[] = [];
    myFollowers.forEach((u: FollowerFromState) => {
      const follower = {
        firstName: u.first_name,
        lastName: u.last_name,
        id: u.user_id,
      };
      l.push(follower);
    });
    setListFollowers(l);
  },[]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const pr = parseInt((event.target as HTMLInputElement).value);
    setValue(pr as Privacy);
  };

  const imgCheck = (image: FileList): boolean => {
    return checkImage(image, setErrors);
  };

  let chosenUsers: Follower[] = [];
  const chosen = (users: readonly Follower[]) => {
    chosenUsers = users as Follower[];
  };

  const newPost = async (data: NewPostForm) => {
    if (!props.parentPrivacy) {
      data.privacy = value as Privacy;
      data.parent_id = 0;
    } else {
      data.privacy = props.parentPrivacy;
      data.parent_id = param;
    }
    let check = true;
    if (chosenUsers.length !== 0) {
      data.userList = chosenUsers.map((user) => user.id);
    }

    if (data.privacy !== 1) {
      let id = getTokenId();
      if (data.privacy == 2) {
        data.userList = listOfFollowers.map((f) => f.id);
      }
      data.userList?.push(id);
    }
    if (data.image.length !== 0) {
      check = imgCheck(data.image);
      data.imgString = (await getBase64(data.image[0])
        .then((str) => {
          return str;
        })
        .catch((e) => alert(e))) as string;
    } else {
      data.imgString = "";
    }

    if (check) {
      try {
        const response = await postService.addNewPost(data);
        handleClose();

        if (!props.parentPrivacy) {
          dispatch(updatePosts(response));
        } else dispatch(updateComments(response));
      } catch (e) {
        console.error(e);
        const errState = {
          text: "Failed to add post",
          severity: "warning",
        };
        dispatch(setAlert(errState));
      }
    } else {
      const errState = {
        text: "Can't upload image (wrong extentsion or image too large) !",
        severity: "warning",
      };
      dispatch(setAlert(errState));
    }
  };

  return (
    <Modal
      open={open}
      onClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
    >
      <Box
        component={"form"}
        onSubmit={handleSubmit(newPost)}
        className={"new_post"}
        sx={{
          boxShadow: 24,
          pt: 2,
          px: 4,
          pb: 3,
        }}
      >
        <div className={"title"}>
          {!props.parentPrivacy && (
            <Input
              placeholder="Title"
              inputProps={ariaLabel}
              fullWidth
              {...register("title")}
              required
            />
          )}
          <Button className={"close"} onClick={handleClose}>
            <CloseIcon />
          </Button>
        </div>
        <div className={"content"}>
          <TextField
            fullWidth
            placeholder="Content"
            multiline
            minRows={5}
            maxRows={10}
            {...register("content")}
            required
          />
        </div>

        <Input sx={{ mb: 3 }} type={"file"} {...register("image")} />
        {!props.parentPrivacy && (
          <div>
            <FormControl id={""} sx={{ ml: 1, mb: 3 }}>
              Who can see this post?
              <RadioGroup
                row
                name="post-privacy"
                value={value}
                onChange={handleChange}
              >
                <FormControlLabel
                  value={Privacy.Public}
                  control={<Radio />}
                  label="All users"
                  onChange={(e) => {
                    e.preventDefault();
                    setFollowers(false);
                    chosenUsers = [];
                  }}
                />
                <FormControlLabel
                  value={Privacy.Private}
                  control={<Radio />}
                  label="Followers" //Friends?
                  onChange={(e) => {
                    e.preventDefault();
                    setFollowers(false);
                    chosenUsers = [];
                  }}
                />
                <FormControlLabel
                  value={Privacy.StrictlyPrivate}
                  control={<Radio />}
                  label="Chosen ones"
                  onChange={(e) => {
                    e.preventDefault();
                    setFollowers(true);
                  }}
                />
              </RadioGroup>
              {followers && (
                <UserList
                  followers={listOfFollowers}
                  sendBack={chosen}
                ></UserList>
              )}
            </FormControl>
          </div>
        )}
        <div>
          <Button sx={{ width: 125 }} type={"submit"} variant="contained">
            Add
          </Button>
        </div>
      </Box>
    </Modal>
  );
}
