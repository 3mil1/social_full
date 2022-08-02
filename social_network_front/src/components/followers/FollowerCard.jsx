import { Avatar, Box, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import "./follower.scss";

const FollowerCard = ({ data }) => {
  let redirect = useNavigate();
  return (
    <Box
      className="follower_card"
      sx={{ backgroundColor: "primary.dark" }}
      onClick={() => {
        redirect(`/profile/${data.user_id}`);
      }}
    >
      <Avatar
        sx={{
          mt: 1,
          bgcolor: "secondary.main",
          width: 50,
          height: 50,
          border: "2px solid white",
        }}
        alt={data.first_name}
        src={data.image}
      ></Avatar>
      <Typography variant="h6" color={"white"} sx={{ mt: -0.5 }}>
        {data.first_name}
      </Typography>
    </Box>
  );
};

export default FollowerCard;
