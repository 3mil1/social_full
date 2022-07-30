import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import "./login.scss";
// Redux
import { useDispatch} from "react-redux";
// Material UI
import {
  Avatar,
  Box,
  Button,
  Container,
  TextField,
  Typography,
} from "@mui/material";
import LoadingButton from "@mui/lab/LoadingButton";
import { PeopleAlt } from "@mui/icons-material";
//  Services
import authService from "../../utilities/user-service";
import ProfileService from "../../utilities/profile_service";
import FollowerService from "../../utilities/follower_service";
import { setAlert } from "../../store/alertSlice";
import WsApi from "../../utilities/ws";
import * as helper from "../../helpers/HelperFuncs";
import GroupService from "../../utilities/group_service";

export default function Login() {
  const profile_service = ProfileService();
  const follower_service = FollowerService();
  const group_service = GroupService();
  const { handleSubmit, register } = useForm<FormInput>();
  const dispatch = useDispatch();
  let redirect = useNavigate();


  interface FormInput {
    email: string;
    password: string;
  }

  const onSubmit = async (data: FormInput) => {
    try {
      await authService.login(data.email, data.password);
      // const response = await authService.login(data.email, data.password);
      // dispatch(update(response));
      // await profile_service.getMyInfo();
      // await follower_service.getMyFollowers();
      profile_service.checkAuth();
      profile_service.getAllUsers()
      group_service.getAllGroups()
      //ws connection
      let id = helper.getTokenId();

      WsApi.start(id, dispatch);

      redirect("/homepage", { replace: true });
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
        const errorState = {
          text: "Mail or password is incorrect",
          severity: "warning",
        };
        dispatch(setAlert(errorState));
      } else {
        console.log(e);
        const errorState = {
          text: "Unknown error occurred",
          severity: "error",
        };
        dispatch(setAlert(errorState));
      }
    }
  };

  // console.log("Here!");
  return (

    <Container component="main" maxWidth="xs" className={"Login"}>
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Avatar
          sx={{
            m: 2,
            mt: 7,
            bgcolor: "secondary.main",
            width: 56,
            height: 56,
          }}
        >
          <PeopleAlt fontSize={"large"} />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <Box
          component={"form"}
          onSubmit={handleSubmit(onSubmit)}
          data-testid={"login-form"}
        >
          <TextField
            required
            label="Mail/Username"
            variant="outlined"
            type={"email"}
            margin="normal"
            {...register("email")}
            data-testid={"email-input"}
          />
          <TextField
            required
            variant="outlined"
            label="Password"
            type={"password"}
            margin="normal"
            {...register("password")}
            data-testid={"pwd-input"}
          />
          <>
            <LoadingButton
              sx={{ mt: 2 }}
              size={"large"}
              variant="contained"
              type={"submit"}
              data-testid={"submit-btn"}
            >
              Login
            </LoadingButton>

            <Button
              sx={{ mt: 2 }}
              size={"large"}
              variant="contained"
              onClick={(e) => {
                e.preventDefault();
                redirect("/register");
              }}
            >
              Register
            </Button>
          </>
        </Box>
      </Box>
    </Container>
  );
}
