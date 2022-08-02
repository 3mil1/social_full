export function parseDate(str: string, postDate: boolean): string {
  const date = new Date(str);
  const min = (date.getMinutes() < 10 ? "0" : "") + date.getMinutes();
  const hour = (date.getHours() < 10 ? "0" : "") + date.getHours();
  const day = (date.getDate() < 10 ? "0" : "") + date.getDate();
  const month = (date.getMonth() + 1 < 10 ? "0" : "") + (date.getMonth() + 1);
  const year = date.getFullYear();
  return postDate
    ? day + "/" + month + "/" + year + " " + hour + ":" + min
    : hour + ":" + min + " " + day + "-" + month + "-" + year;
}
