import { Box, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import CoPresentIcon from "@mui/icons-material/CoPresent";
import "./group.scss";

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
            sx={{ mt: -0.5, letterSpacing: 2 }}
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
