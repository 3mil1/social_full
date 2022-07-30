
const  todayTime = () => {
    let currentTime = new Date();
    let hour = currentTime.getHours();
    let minutes = currentTime.getMinutes();
    return  `${hour}:${minutes < 10 ? "0" + minutes : minutes}`
}

const todayDate = () => {
    let currentTime = new Date();
    let month = currentTime.getMonth() + 1;
    month = month < 10 ? '0' + month : month;
    let day = currentTime.getDate();
    day = day < 10 ? '0' + day : day;
    let year = currentTime.getFullYear();
    var todaysDate = year + '-' + month + '-' + day;
    return todaysDate;
}


const isFuture = (current, target) => {
    let currentArr = current.split('-').map(Number);
    let targetArr = target.split('-').map(Number);
    if (currentArr[0] < targetArr[0]) return true;
    if (currentArr[0] == targetArr[0] && currentArr[1] < targetArr[1])
    return true;
    if (
    currentArr[0] == targetArr[0] &&
    currentArr[1] == targetArr[1] &&
    currentArr[2] < targetArr[2]
    )
    return true;
    return false;
}

const  calcTime = time => {
    let arr = time.split(':').map(Number);
    return Number(arr[0] * 60 + arr[1]);
};

export { todayTime, todayDate, isFuture, calcTime };