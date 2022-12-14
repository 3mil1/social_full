import React from "react";
import Snackbar from "@mui/material/Snackbar";
import { setAlert } from "../store/alertSlice";
import { Alert } from "@mui/material";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../store/store";

export function AlertSnackbar() {
  const dispatch = useDispatch();
  const alertText = useSelector((state: RootState) => state.alert.text);
  const isOpen = useSelector((state: RootState) => state.alert.isOpen);
  const severity = useSelector((state: RootState) => state.alert.severity);

  const handleClose = (
    event?: React.SyntheticEvent | Event,
    reason?: string
  ) => {
    const errorState = {
      text: "",
      severity: undefined,
    };
    if (reason === "clickaway") {
      dispatch(setAlert(errorState));
    }
    dispatch(setAlert(errorState));
  };

  if (!isOpen) {
    return <></>;
  }

  return (
    <Snackbar
      anchorOrigin={{ vertical: "top", horizontal: "center" }}
      style={{ top: "25%" }}
      open={isOpen}
      autoHideDuration={6000}
      onClose={handleClose}
    >
      <Alert
        onClose={handleClose}
        severity={severity}
        elevation={6}
        variant="filled"
      >
        {alertText}
      </Alert>
    </Snackbar>
  );
}
