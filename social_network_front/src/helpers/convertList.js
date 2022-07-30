// StoreLists themselves are read only (copy that store list wont help, cuz configurations wont change)

const convertListToMutable = arr => {
    let convertedList = [];
    arr.forEach(obj => {
        convertedList.push(JSON.stringify(obj))
    });

    return convertedList.map(obj => JSON.parse(obj));
};

export default convertListToMutable;
