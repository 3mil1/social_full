const checkError = err => {
  if (err.response) {
    console.log(err.message);
    console.log(err.status);
    console.log(err.header);
  } else {
    console.log('Different Error ', err);
  }
};

export default checkError