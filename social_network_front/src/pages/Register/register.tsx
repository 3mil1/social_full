import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import "./register.scss";
// UI material
import { DatePicker } from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import {
  Avatar,
  Box,
  Button,
  Container,
  Grid,
  TextField,
  Typography,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
// other
import * as helper from "../../helpers/HelperFuncs";
import userService from "../../utilities/user-service";
import { format } from "date-fns";

export interface RegisterForm {
  first_name: string;
  last_name: string;
  dob: string;
  email: string;
  password: string;
  repeat_password: string;
  nickname: string;
  image_path: string;
  desc: string;
}

export default function Register() {
  const [value, setValue] = useState<Date | null>(null);
  let redirect = useNavigate();
  const { handleSubmit, register } = useForm<RegisterForm>();
  const [errors, setErrors] = useState([]);
  const myMaxDate = new Date(
    new Date().getFullYear() - 18,
    new Date().getMonth(),
    new Date().getDate()
  );
  const onSubmit = async (user: RegisterForm) => {
    let flag = true;
    setErrors([]);
    user.nickname = user.nickname ? user.nickname : "";
    user.image_path = user.image_path ? user.image_path : "";
    user.desc = user.desc ? user.desc : "";

    // @ts-ignore
    if (user.password !== user.repeat_password) {
      flag = false;
      // @ts-ignore
      setErrors((oldArray) => [...oldArray, "* Passwords do not match"]);
    }
    // @ts-ignore
    if (user.image_path.length !== 0 && user.image_path[0].name !== "") {
      if (!helper.checkImage(user.image_path, setErrors)) flag = false;
      try {
        // @ts-ignore
        user.image_path = await helper
          .getBase64(user.image_path[0])
          .then((base64) => {
            return base64;
          });
      } catch (e) {
        // @ts-ignore
        setErrors((oldArray) => [...oldArray, "ERROR WITH IMAGE UPLOAD"]);
      }
    }

    if (flag) {
      try {
        if (user.image_path.length == 0) user.image_path = "";
        // @ts-ignore
        const formatted = format(value, "dd-MM-yyyy");
        user.dob = formatted;
        await userService.register(user);
        redirect("/");
      } catch (e) {
        if (e instanceof Error) {
          console.log("1", e);
          // @ts-ignore
          setErrors((oldArray) => [
            ...oldArray,
            // @ts-ignore
            `${helper.capitalize(e.response?.data.message || e.message)}`,
          ]);
        } else {
          console.error(e);
          redirect("/");
        }
      }
    }
  };

  // @ts-ignore
  return (
    <Container component="main" maxWidth="md" className={"Register"}>
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Button
          variant="contained"
          className="back_btn"
          onClick={(e) => {
            e.preventDefault();
            redirect("/");
          }}
        >
          <CloseIcon />
        </Button>
        <Avatar sx={{ m: 1, mt: 4, bgcolor: "secondary.main" }}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Register
        </Typography>

        {errors && (
          <div className="errors">
            {errors.map((err, i) => (
              <div key={i}>{err}</div>
            ))}
          </div>
        )}

        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 2 }}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                sx={{ m: 1, width: "280px" }}
                required
                inputProps={{ minLength: 2, maxLength: 30 }}
                label="First Name"
                type="text"
                margin="dense"
                variant="standard"
                {...register("first_name")}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                inputProps={{ minLength: 2, maxLength: 30 }}
                label="Last Name"
                type="text"
                margin="dense"
                variant="standard"
                sx={{ m: 1, width: "280px" }}
                {...register("last_name")}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                required
                inputProps={{ minLength: 2, maxLength: 30 }}
                label="Email"
                type="email"
                variant="standard"
                sx={{ m: 1, width: "280px" }}
                {...register("email")}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                inputProps={{ minLength: 6, maxLength: 30 }}
                label="Password"
                type="password"
                margin="dense"
                variant="standard"
                sx={{ m: 1, width: "280px" }}
                {...register("password")}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                inputProps={{ minLength: 6, maxLength: 30 }}
                label="Repeat password"
                type="password"
                margin="dense"
                variant="standard"
                sx={{ m: 1, width: "280px" }}
                {...register("repeat_password")}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker
                  label="Birthday(18+)"
                  {...register("dob")}
                  onChange={(value) => setValue(value)}
                  value={value}
                  minDate={new Date("1922-01-01")}
                  maxDate={myMaxDate}
                  renderInput={(params) => <TextField {...params} />}
                  inputFormat={"dd-MM-yyyy"}
                />
              </LocalizationProvider>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Nickname"
                inputProps={{ maxLength: 10 }}
                type="text"
                margin="dense"
                variant="standard"
                sx={{ m: 1, width: "280px" }}
                {...register("nickname")}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Avatar"
                type="file"
                margin="normal"
                variant="standard"
                {...register("image_path")}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                inputProps={{ maxLength: 30 }}
                label="About Me"
                type="text"
                multiline
                rows={3}
                sx={{ mt: 2, width: "316px" }}
                {...register("desc")}
              />
            </Grid>
            <Grid item xs={12} sm={6} spacing={6} />
            <Grid item xs={12} sm={6} spacing={6}>
              <Button
                variant="contained"
                type={"submit"}
                sx={{ mt: 3 }}
                size={"large"}
              >
                Sign Up
              </Button>
            </Grid>
          </Grid>
        </Box>
      </Box>
    </Container>
  );
}
