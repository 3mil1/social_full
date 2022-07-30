// @ts-ignore
export const checkImage = (imageList, func): boolean => {
  console.log("CHECKIMAGE ---> ", imageList);
  
  let regex = /[a-zA-Z\d](.jpg|.JPG|.gif|.GIF|.jpeg|.JPEG|.png|.PNG)/gm;
  let save = true;
  if (imageList.length != 0) {
    // @ts-ignore
    if (Number(imageList[0].size) >= 4194304) {
      //4mb
      console.log("FRONT END - %cImage size is Large", "color:red");
      // @ts-ignore
      func((oldArray) => [...oldArray, "* Image size is too large"]);
      save = false;
    }
    // @ts-ignore
    if (!regex.test(imageList[0].name)) {
      console.log("FRONT END - %cImage file extension is wrong", "color:red");
      // @ts-ignore
      func((oldArray) => [
        ...oldArray,
        "* Only file extension with (jpeg,jpg,png,gif) are allowed",
      ]);
      save = false;
    }
  }
  return save;
};

// @ts-ignore
export const getBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
    reader.readAsDataURL(file);
  });
};


// @ts-ignore
export const convertImg = async (image,func) => { 
    if(image.length !== 0) {
      if(checkImage(image,func)){
        func([])
        const resp = await getBase64(image[0]).then(base64 => base64)
        return resp
      }
    }
  }