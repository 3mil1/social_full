import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ProfileService from "../utilities/profile_service";
import Follow_btn from "./buttons/follower_btn";
import * as helper from "../helpers/HelperFuncs";
//  MUI Material
import { Avatar, Button, Input, Typography } from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import EditIcon from "@mui/icons-material/Edit";
import RemoveIcon from "@mui/icons-material/Remove";
// Redux
import { useSelector } from "react-redux";
import "./styles/profile_info.scss";

const ProfileInfo = () => {
  const profile_service = ProfileService();
  const storeInfo = useSelector((state) => state);
  const [updateing, setUpdateing] = useState(false);
  const [isPrivate, setIsPrivate] = useState(null);
  const [img, setImg] = useState(null);
  const [errors, setErrors] = useState([]);
  const [data, setDatas] = useState({});
  const [myProfile, setMyProfile] = useState(false);
  let { id } = useParams();
  let updateInfo = {
    nickname: data.nickname,
    about_me: data.about_me,
    user_img: data.user_img,
    is_private: isPrivate,
  };

  const convertImg = async (image) => {
    if (image.length !== 0) {
      if (helper.checkImage(image, setErrors)) {
        setErrors([]);
        const resp = await helper.getBase64(image[0]).then((base64) => base64);
        return resp;
      }
    }
  };

  const followingAlready = () => {
    return storeInfo.followers.followers.filter(user => user.user_id === id).length != 0;
  }

  const handleUpdate = async (info) => {
    if (errors.length == 0) {
      info.user_img = img;
      profile_service.updateProfileInfo(info);
    }
  };

  const updateData = () => {
    setIsPrivate(storeInfo.profile.info.is_private);
    if (id == "me") {
      setMyProfile(true);
      setDatas(storeInfo.profile.info);
      setImg(data.user_img);
    } else {
      setUpdateing(false);
      setMyProfile(false);
      profile_service.getUserInfo(id).then((res) => {
        setIsPrivate(res.is_private);
        setDatas({ ...res, id });
      });
    }
  };

  useEffect(() => {
    followingAlready()
    updateData();
    if (!updateing && myProfile) {
      setImg(data.user_img);
      setErrors([]);
    }
  }, [id, myProfile, updateing]);

  return (
    <div className="user_info_container">
      {errors && (
        <div className="errors">
          {errors.map((err, i) => (
            <div key={i}>{err}</div>
          ))}
        </div>
      )}
      <div className="left_side">
        {myProfile && (
          <div
            className="setting_btn"
            onClick={() => {
              setUpdateing(!updateing);
            }}
          >
            <SettingsIcon className="gear" />
          </div>
        )}

        {/* IF User is updating settings */}
        {updateing && myProfile && (
          <>
            <div className="flex privacy_btn_wrapper ">
              Public
              <div
                onClick={() => {
                  setIsPrivate(!isPrivate);
                  updateInfo.is_private = isPrivate;
                }}
                className={isPrivate ? "public private " : "public"}
              ></div>{" "}
              Private
            </div>

            <div className="profile_image_container">
              <Avatar
                sx={{ width: 80, height: 80, opacity: 0.8 }}
                alt={data.first_name}
                src={img}
              />
              <label className="flex" id="edit_icon" htmlFor="avatar">
                <EditIcon />
              </label>
              <input
                id="avatar"
                type="file"
                accept="image/*,.png, .jpg, .jpeg, .gif"
                onChange={() => {
                  convertImg(document.getElementById("avatar").files).then(
                    (res) => setImg(res)
                  );
                }}
              />
              <label
                className="flex"
                id="remove_icon"
                onClick={() => {
                  setImg("");
                  updateInfo.user_img = "";
                }}
              >
                <RemoveIcon />
              </label>
            </div>

            <Input
              className="update_field"
              type={"text"}
              defaultValue={updateInfo.nickname}
              onInput={(e) => {
                updateInfo.nickname = e.target.value;
              }}
            >
              {data.nickname}
            </Input>
            <Button
              onClick={() => {
                handleUpdate(updateInfo);
              }}
            >
              Update
            </Button>
          </>
        )}

        {!updateing && (
          <>
            <Typography variant="h6" gutterBottom>
              Profile is {data.is_private ? "private" : "public"}
            </Typography>
            <Avatar
              sx={{ width: 100, height: 100 }}
              alt={data.first_name}
              src={data.user_img}
            />
            <Typography sx={{ margin: "1em", fontSize: "20px" }}>
              {data.nickname}
            </Typography>

            {!myProfile && <Follow_btn isPrivate={isPrivate} />}
          </>
        )}
      </div>

      <div className="right_side">
        <p> First Name : {data.first_name}</p>
        <p> Last Name : {data.last_name}</p>

        {!myProfile && (!isPrivate || followingAlready()) && (
          <>
            <p> Email : {data.email}</p>
            <p> Birthday : {data.birth_day}</p>
            <p> About me : {data.about_me}</p>
          </>
        )}

        {myProfile && (
          <>
            <p> Email : {data.email}</p>
            <p> Birthday : {data.birth_day}</p>
            {!updateing ? (
              <p> About me : {data.about_me}</p>
            ) : (
              <div className="flex about_me">
                <p> About me : </p>
                <Input
                  type={"text"}
                  defaultValue={updateInfo.about_me}
                  onInput={(e) => {
                    updateInfo.about_me = e.target.value;
                  }}
                ></Input>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default ProfileInfo;
