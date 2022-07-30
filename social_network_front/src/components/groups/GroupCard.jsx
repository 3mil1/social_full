import { Box, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import CoPresentIcon from "@mui/icons-material/CoPresent";
import "./group.scss";
import { useState } from "react";

/* 
  type GroupReply struct{
      Id                  int     `json:"id"`
      Title               string  `json:"title"`
      Description         string  `json:"description"`
      CreatorId           string  `json:"creator_id"`
      CreatorFirstName    string  `json:"creator_first_name"`
      CreatorLastName     string  `json:"creator_last_name"`
      Members             int     `json:"members"`
  }
*/

// const group = {
//     id : "dev",
//     title : "kmds",
//     description : "oke",
//     creator_id : "12-234",
//     creator_first_name : "Sil",
//     creator_last_name : "ver",
//     members : 10,
// }

const GroupCard = ({ data, myInfo }) => {
  let redirect = useNavigate();
  return (
    <Box
      className="group_card"
      sx={{ backgroundColor: "primary.dark" }}
      onClick={() => {
        redirect(`/group/${data.id}`);
      }}
    >
      <div className="header">
        {myInfo ? (
          <Typography
            variant="h5"
            color={"white"}
            sx={{ mt: -0.5,letterSpacing : 2 }}
          >
            {data.title}
          </Typography>
        ) : (
          <div className="title">
            <Typography
              variant="h5"
              color={"white"}
              sx={{ mt: -0.5, textDecoration: "underline" }}
            >
              {data.title}
            </Typography>
            <Typography
              variant="h5"
              color={"white"}
              sx={{ m: "4px", fontSize: "10px", opacity: 0.8 }}
            >
              ( {data.creator_first_name} {data.creator_last_name} )
            </Typography>
          </div>
        )}
        <Typography
          variant="h6"
          color={"white"}
          sx={{
            mt: -0.5,
            position: "relative",
            display: "flex",
            placeItems: "center",
            width: "40px",
            height: "40px",
          }}
        >
          <CoPresentIcon sx={{ fontSize: "30px", zIndex: 0, opacity: 0.4 }} />
          <Typography
            sx={{
              position: "absolute",
              right: "0px",
              top: "0px",
              zIndex: 2,
              fontSize: "20px",
            }}
          >
            {data.members}
          </Typography>
        </Typography>
      </div>
      <Typography variant="p" color={"white"} sx={{ mt: -0.5 }}>
        " {data.description} "
      </Typography>
    </Box>
  );
};

export default GroupCard;
